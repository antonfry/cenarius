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

func (a *agent) Shutdown() {

}

// r, w := io.Pipe()
// m := multipart.NewWriter(w)
// go func() {
//     defer w.Close()
//     defer m.Close()
//     part, err := m.CreateFormFile("myFile", "foo.txt")
//     if err != nil {
//         return
//     }
//     file, err := os.Open(name)
//     if err != nil {
//         return
//     }
//     defer file.Close()
//     if _, err = io.Copy(part, file); err != nil {
//         return
//     }
// }()
// http.Post(url, m.FormDataContentType(), r)
