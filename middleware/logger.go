package middleware

import (
	"log"
	"net/http"
	"time"
)

// RequestLogger returns a handler function that will log info about
// the request. Then it calls the provided handler function.
func RequestLogger(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	// Get the current time.
	start := time.Now()

	// Once the handler call proceeding this defer
	// is complete, log how long the request took.
	defer func() {
		d := time.Since(start)
		log.Printf("%s %s %s", req.Method, req.URL.Path, d)
	}()

	// Call the handler that was provided.
	next(res, req)
}
