package agent

import (
	"bytes"
	"cenarius/internal/model"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type agent struct {
	client http.Client
	config *Config
	logger *logrus.Logger
}

// NewServer returns new server object
func NewAgent(config *Config) *agent {
	a := &agent{
		client: http.Client{},
		config: config,
		logger: logrus.New(),
	}
	return a
}

// Start starts the agent
func (a *agent) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	a.configureLogger()
	//a.register(ctx)
	a.getLogingWithPassword(ctx)
	return nil
}

// Stop stops the agent
func (a *agent) Shutdown() {

}

// configureLogger configures logger
func (a *agent) configureLogger() error {
	level, err := logrus.ParseLevel(a.config.LogLevel)
	if err != nil {
		return err
	}
	a.logger.SetLevel(level)
	return nil
}

// write2Buffer writes jsonData to buf
func (a *agent) write2Buffer(jsonData []byte, buf *bytes.Buffer) {
	if a.config.Encode {
		gzipData := gzip.NewWriter(buf)
		if _, err := gzipData.Write(jsonData); err != nil {
			a.logger.Errorf("write2Buffer err: %s", err.Error())
			return
		}
		if err := gzipData.Close(); err != nil {
			a.logger.Errorf("write2Buffer err: %s", err.Error())
			return
		}
	} else {
		_, _ = io.WriteString(buf, string(jsonData))
	}
}

func (a *agent) sendRequest(ctx context.Context, path string, method string, v any) {
	endpoint := fmt.Sprintf("http://%v/%s", a.config.Host, path)
	var buf bytes.Buffer
	a.logger.Debug("sendRequest endpoint: ", endpoint)
	if v != nil {
		jsonData, err := json.Marshal(v)
		if err != nil {
			a.logger.Errorf("sendRequest err: %s", err.Error())
			return
		}
		a.logger.Infof("sendRequest jsonData: %s", string(jsonData))
		a.write2Buffer(jsonData, &buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, &buf)
	if err != nil {
		a.logger.Errorf("sendRequest err: %s", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if a.config.Encode {
		req.Header.Set("Content-Encoding", "gzip")
	}
	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Errorf("sendRequest err: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	a.logger.Infof("sendRequest response code: %d", resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		a.logger.Errorf("sendRequest err: %s", err.Error())
	}
	fmt.Printf("Response body: %v", bytes.NewBuffer(bodyBytes).String())
}

func (a *agent) register(ctx context.Context) {
	uri := "api/v1/user/register"
	m := &model.User{Login: a.config.Login, Password: a.config.Password}
	a.sendRequest(ctx, uri, http.MethodPost, m)
}

func (a *agent) login(ctx context.Context) {
	uri := "api/v1/user/register"
	m := &model.User{Login: a.config.Login, Password: a.config.Password}
	a.sendRequest(ctx, uri, http.MethodPost, m)
}

func (a *agent) getLogingWithPassword(ctx context.Context) {
	uri := "api/v1/private/loginwithpasswords"
	a.sendRequest(ctx, uri, http.MethodGet, nil)
}
