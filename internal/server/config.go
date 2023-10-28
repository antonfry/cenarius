package server

type Config struct {
	Bind           string `json:"bind"`
	Mode           string `json:"mode"`
	LogLevel       string `json:"log_level"`
	DatabaseDsn    string `json:"database_url"`
	SessionKey     string `json:"session_key"`
	SecretFilePath string `json:"secret_file_path"`
}

func NewConfig() *Config {
	return &Config{
		Bind:           ":8080",
		Mode:           "HTTP",
		LogLevel:       "INFO",
		DatabaseDsn:    "postgres://localhost/cenarius?sslmode=disable",
		SessionKey:     "cenarius",
		SecretFilePath: "/tmp/cenarius",
	}
}
