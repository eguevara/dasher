package config

// AppConfig stores application settings.
type AppConfig struct {
	Debug              bool         `json:"debug"`
	Environment        string       `json:"environment"`
	PrivateKeyFilename string       `json:"pemFilePath"`
	Version            string       `json:"version"`
	Address            string       `json:"address"`
	ShutdownTimeout    int64        `json:"shutdownTimeout"`
	AnalyticsOAuth     *OAuthConfig `json:"analyticsOAuth"`
}

// OAuthConfig stores oauth settings
type OAuthConfig struct {
	PrivateFilePath string   `json:"pemFilePath"`
	ServiceEmail    string   `json:"serviceEmail"`
	Scopes          []string `json:"scopes"`
}
