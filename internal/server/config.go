package server

type Config struct {
	Bind          string `json:"bind"`
	Mode          string `json:"mode"`
	LogLevel      string `json:"log_level"`
	DatabaseDsn   string `json:"database_url"`
	TrustedSubnet string `json:"trusted_subnet"`
}

func NewConfig() *Config {
	return &Config{
		Bind:        ":8080",
		Mode:        "HTTP",
		LogLevel:    "INFO",
		DatabaseDsn: "",
	}
}
