package server

type Config struct {
	Bind           string `json:"bind" toml:"bind,omitempty"`
	LogLevel       string `json:"log_level" toml:"log_level,omitempty"`
	DatabaseDsn    string `json:"database_url" toml:"database_url,omitempty"`
	SessionKey     string `json:"session_key" toml:"session_key,omitempty"`
	SecretFilePath string `json:"secret_file_path" toml:"secret_file_path,omitempty"`
}

func NewConfig() *Config {
	return &Config{
		Bind:           ":8080",
		LogLevel:       "INFO",
		DatabaseDsn:    "postgres://localhost/cenarius?sslmode=disable",
		SessionKey:     "cenarius",
		SecretFilePath: "/tmp/cenarius",
	}
}
