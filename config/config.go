package config

// AppConfig stores application settings.
type AppConfig struct {
	Debug              bool          `json:"debug"`
	Environment        string        `json:"environment"`
	PrivateKeyFilename string        `json:"pemFilePath"`
	Version            string        `json:"version"`
	Address            string        `json:"address"`
	ShutdownTimeout    int64         `json:"shutdownTimeout"`
	AnalyticsOAuth     *OAuthConfig  `json:"analyticsOAuth"`
	BooksOAuth         *OAuthConfig  `json:"booksOAuth"`
	BooksShelf         *string       `json:"booksShelf"`
	BooksVolumesMax    *int          `json:"booksVolumesMax"`
	BooksVolumesFields *string       `json:"booksVolumesFields"`
	GitHub             *GitHubConfig `json:"github"`
	CatalogFilename    string        `json:"catalogFilePath"`
}

// OAuthConfig stores oauth settings
type OAuthConfig struct {
	PrivateFilePath  string   `json:"pemFilePath"`
	ServiceEmail     string   `json:"serviceEmail"`
	Scopes           []string `json:"scopes"`
	ImpersonateEmail *string  `json:"impersonateEmail"`
}

// GitHubConfig stores any of the github options to connect.
type GitHubConfig struct {
	Token string `json:"token"`
	Repo  string `json:"repo"`
	Owner string `json:"owner"`
}
