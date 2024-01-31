# Compile
`make`

## Linux
`make build_linux`
## Windows
`make build_windows`
## MacOS
`make build_macos`

# Tests
`make tests`


# How to Run

## Docker compose
Server
`make compose`
Agent:
`./cmd/cenarius/cenarius -m agent -host localhost:443`


## Server mode
`./cmd/cenarius/cenarius -m server`

## Agent mode 
`./cmd/cenarius/cenarius -m agent`

# Configuration

## Command line flags
```  
  -conf string
    	path to toml conf (default "conf/conf.toml")
  -databaseDSN string
    	Database DNS for server
  -host string
    	Server address
  -logLevel string
    	LogLevel
  -login string
    	Login for agent
  -m string
    	server or agent
  -password string
    	Password for agent
  -secretFilePath string
    	Storage path for secret files
```

## Environment variables
### server
```CENARIUS_LOG_LEVEL - logging level
CENARIUS_SERVER_BIND - Address to bind server
CENARIUS_DATABASEDSN - Postgre dsn(Example: "postgres://postgres:password@localhost:5432/cenarius_test?sslmode=disable",)
CENARIUS_SECRET_STORAGE_PATH - Path to storage for secret files```
### agent
```CENARIUS_LOG_LEVEL - logging level
CENARIUS_SERVER_ADDR - cenarius server address
CENARIUS_LOGIN - cenarius server login
CENARIUS_PASSWORD - cenarius server password```
