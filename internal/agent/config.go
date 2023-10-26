package agent

type Config struct {
	Host     string `json:"host"`
	Mode     string `json:"mode"`
	LogLevel string `json:"log_level"`
	Secret   string `json:"secret"`
	GZip     bool   `json:"gzip"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewConfig() *Config {
	return &Config{
		Host:     "localhost:8080",
		Mode:     "HTTP",
		LogLevel: "INFO",
		Secret:   "",
		GZip:     false,
		Login:    "AgentUser",
		Password: "AgentPassword",
	}
}
