package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"

	"github.com/eguevara/dasher/api"
	"github.com/eguevara/dasher/config"
	"github.com/eguevara/dasher/middleware"
)

const (
	version        = "1.0.0"
	realtimePrefix = "/v1/realtime"
	metricsPath    = "/metrics"
	healthPath     = "/health"
	versionPath    = "/version"
)

// init is called before main. We are using init to customize logging output.
func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func main() {

	// Set up application configuration.
	cfg := appConfig()

	// Create a new server and set timeout values.
	server := http.Server{
		Addr:           cfg.Address,
		Handler:        appHandlers(cfg),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start the listener.
	go func() {
		log.Printf("start: Listening on %s\n", cfg.Address)
		log.Println("start: Process ID", os.Getpid())

		if err := server.ListenAndServe(); err != nil {
			log.Println("ListenAndServe returns an error", err)
		}
	}()

	// Listen for an interrupt signal from the OS.
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)

	// Wait for a signal to shutdown.
	log.Printf("shutdown: Signal %v", <-signalChan)

	// Create a context to attempt a graceful 5 second shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeout)*time.Second)
	defer cancel()

	// Attempt the graceful shutdown by closing the listener and
	// completing all inflight requests.
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown : Graceful shutdown did not complete in %v : %v", time.Duration(cfg.ShutdownTimeout)*time.Second, err)

		// Looks like we timedout on the graceful shutdown. Kill it hard.
		if err := server.Close(); err != nil {
			log.Printf("shutdown : Error killing server : %v", err)
		}
	}

	log.Println("shutdown: Graceful complete")

}

func appHandlers(cfg *config.AppConfig) http.Handler {

	// Using gorilla mux for richer routing
	r := mux.NewRouter()
	r.Handle(versionPath, api.VersionHandler(version)).Methods("GET")

	books := api.NewBooksHandler(cfg)
	r.HandleFunc("/v1/books", books.List).Methods("GET")
	r.HandleFunc("/v1/books/{id}", books.ListNotesByVolumeID).Methods("GET")

	r.HandleFunc(healthPath, api.HealthHandler).Methods("GET")
	r.Handle(metricsPath, promhttp.Handler())

	n := negroni.New()

	// Add Request Logger Middleware
	n.UseFunc(middleware.RequestLogger)

	// Add Application handler Middleware
	n.UseFunc(middleware.AppHandler)

	// Apply the mux to negroni.
	n.UseHandler(r)

	return n
}

func appConfig() *config.AppConfig {

	var (
		flagHTTPAddr   = flag.String("web.listen-address", "", "HTTP service address (0.0.0.0:3000)")
		flagConfigPath = flag.String("config.file", "./config.json", "application configuration file")
	)

	flag.Parse()

	configFile := filepath.Join(*flagConfigPath)
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("DASHER error reading config.json: %v", err)
	}

	config := &config.AppConfig{}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("DASHER setting up app configuration: %v", err)
	}

	// Allow address flag override.
	if *flagHTTPAddr != "" {
		config.Address = *flagHTTPAddr
	}

	return config
}
