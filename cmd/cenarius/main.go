package main

import (
	"cenarius/internal/agent"
	"cenarius/internal/server"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type flags struct {
	mode           string
	conf           string
	logLevel       string
	host           string
	databaseDSN    string
	secretFilePath string
	login          string
	password       string
}

type cenariusWorker interface {
	Start() error
	Shutdown()
}

var (
	buildVersion string
	buildDate    string
	buildCommit  string
	flagsData    flags
)

func getServerConfig(conf *server.Config) *server.Config {
	_, err := toml.DecodeFile(flagsData.conf, conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}

func getServerFlags(conf *server.Config) *server.Config {
	if flagsData.logLevel != "" {
		conf.LogLevel = flagsData.logLevel
	}
	if flagsData.host != "" {
		conf.Bind = flagsData.host
	}
	if flagsData.databaseDSN != "" {
		conf.DatabaseDsn = flagsData.databaseDSN
	}
	if flagsData.secretFilePath != "" {
		conf.SecretFilePath = flagsData.secretFilePath
	}
	return conf
}

func getServerEnv(conf *server.Config) *server.Config {
	loglevel, ok := os.LookupEnv("CENARIUS_LOG_LEVEL")
	if ok {
		conf.LogLevel = loglevel
	}
	bind, ok := os.LookupEnv("CENARIUS_SERVER_BIND")
	if ok {
		conf.Bind = bind
	}
	dbDSN, ok := os.LookupEnv("CENARIUS_DATABASEDSN")
	if ok {
		conf.DatabaseDsn = dbDSN
	}
	secretPath, ok := os.LookupEnv("CENARIUS_SECRET_STORAGE_PATH")
	if ok {
		conf.SecretFilePath = secretPath
	}
	return conf
}

func getAgentConfig(conf *agent.Config) *agent.Config {
	_, err := toml.DecodeFile(flagsData.conf, conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}

func getAgentFlags(conf *agent.Config) *agent.Config {
	if flagsData.logLevel != "" {
		conf.LogLevel = flagsData.logLevel
	}
	if flagsData.host != "" {
		conf.Host = flagsData.host
	}
	if flagsData.login != "" {
		conf.Login = flagsData.login
	}
	if flagsData.password != "" {
		conf.Password = flagsData.password
	}
	return conf
}

func getAgentEnv(conf *agent.Config) *agent.Config {

	loglevel, ok := os.LookupEnv("CENARIUS_LOG_LEVEL")
	if ok {
		conf.LogLevel = loglevel
	}
	host, ok := os.LookupEnv("CENARIUS_SERVER_ADDR")
	if ok {
		conf.Host = host
	}
	login, ok := os.LookupEnv("CENARIUS_LOGIN")
	if ok {
		conf.Login = login
	}
	password, ok := os.LookupEnv("CENARIUS_PASSWORD")
	if ok {
		conf.Password = password
	}
	return conf
}

func main() {
	flag.StringVar(&flagsData.mode, "m", "", "server or agent")
	flag.StringVar(&flagsData.conf, "conf", "conf/conf.toml", "path to toml conf")
	flag.StringVar(&flagsData.logLevel, "logLevel", "", "LogLevel")
	flag.StringVar(&flagsData.host, "host", "", "Server address")
	flag.StringVar(&flagsData.databaseDSN, "databaseDSN", "", "Database DNS for server")
	flag.StringVar(&flagsData.secretFilePath, "secretFilePath", "", "Storage path for secret files")
	flag.StringVar(&flagsData.login, "login", "", "Login for agent")
	flag.StringVar(&flagsData.login, "password", "", "Password for agent")
	flag.Parse()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	var worker cenariusWorker
	switch flagsData.mode {
	case "server":
		conf := getServerConfig(server.NewConfig())
		log.Debugf("Conf after file configuration: %v", conf)
		conf = getServerFlags(conf)
		log.Debugf("Conf after flags: %v", conf)
		conf = getServerEnv(conf)
		log.Debugf("Conf after env variables: %v", conf)
		worker = server.NewServer(conf)
	case "agent":
		conf := getAgentConfig(agent.NewConfig())
		log.Debugf("Conf after file configuration: %v", conf)
		conf = getAgentFlags(conf)
		log.Debugf("Conf after flags: %v", conf)
		conf = getAgentEnv(conf)
		log.Debugf("Conf after env variables: %v", conf)
		worker = agent.NewAgent(conf)
	default:
		flag.Usage()
		log.Fatalf("Unknown mode %v", flagsData.mode)
	}

	go func() {
		sig := <-sigs
		log.Infof("OS SIGNAL: %v", sig)
		worker.Shutdown()
	}()
	log.Infof("Build version: %v", buildVersion)
	log.Infof("Build date: %v", buildDate)
	log.Infof("Build commit: %v", buildCommit)
	if err := worker.Start(); err != nil {
		log.Error(err)
	}
}
