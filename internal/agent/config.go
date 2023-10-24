package agent

type Config struct {
	Host     string `json:"host"`
	Mode     string `json:"mode"`
	LogLevel string `json:"log_level"`
	Secret   string `json:"secret"`
	Encode   bool   `json:"encode"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Action   string `json:"action"`
}

func NewConfig() *Config {
	return &Config{
		Host:     "localhost:8080",
		Mode:     "HTTP",
		LogLevel: "INFO",
		Secret:   "",
		Encode:   true,
		Login:    "AgentUser",
		Password: "AgentPassword",
	}
}
