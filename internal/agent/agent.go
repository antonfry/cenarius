package agent

import (
	"github.com/sirupsen/logrus"
)

type agent struct {
	config *Config
	logger *logrus.Logger
}

// NewServer returns new server object
func NewAgent(config *Config) *agent {
	a := &agent{
		config: config,
		logger: logrus.New(),
	}
	return a
}

func (a *agent) Start() {

}
