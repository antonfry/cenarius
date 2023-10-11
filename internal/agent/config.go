package agent

type Config struct {
	ServerURL string `json:"server_url"`
	Mode      string `json:"mode"`
	LogLevel  string `json:"log_level"`
	Secret    string `json:"secret"`
}

func NewConfig() *Config {
	return &Config{
		ServerURL: "localhost:8080",
		Mode:      "HTTP",
		LogLevel:  "INFO",
		Secret:    "",
	}
}
