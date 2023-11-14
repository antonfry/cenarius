package agent

type Config struct {
	Host      string `json:"host" toml:"host,omitempty"`
	LogLevel  string `json:"log_level" toml:"log_level,omitempty"`
	GZip      bool   `json:"gzip" toml:"gzip,omitempty"`
	Login     string `json:"login" toml:"login,omitempty"`
	Password  string `json:"password" toml:"password,omitempty"`
	CacheFile string `json:"cache_file"`
	SecretKey string `json:"secret_key"`
	SecretIV  string `json:"secret_iv"`
}

func NewConfig() *Config {
	return &Config{
		Host:      "localhost:8080",
		LogLevel:  "INFO",
		GZip:      false,
		Login:     "AgentUser",
		Password:  "AgentPassword",
		CacheFile: "/tmp/cenarius.cache",
	}
}
