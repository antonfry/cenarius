package agent

import (
	"bytes"
	"cenarius/internal/model"
	"cenarius/internal/server"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"

	"github.com/sirupsen/logrus"
)

type agent struct {
	client http.Client
	config *Config
	logger *logrus.Logger
}

// NewServer returns new server object
func NewAgent(config *Config) *agent {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal("unable to inizialize cookiejar")
	}
	a := &agent{
		client: http.Client{Jar: jar},
		config: config,
		logger: logrus.New(),
	}
	return a
}

// Start starts the agent
func (a *agent) Start() error {
	ctx := context.Background()
	a.configureLogger()
	a.health(ctx)
	a.health(ctx)
	a.userInput()
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
func (a *agent) write2Buffer(buf *bytes.Buffer, v any) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		a.logger.Errorf("sendRequest err: %s", err.Error())
		return
	}
	a.logger.Infof("sendRequest jsonData: %s", string(jsonData))
	if a.config.GZip {
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

func (a *agent) encodeAuth() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s %s", a.config.Login, a.config.Password)))
}

func (a *agent) getRequest(ctx context.Context, method string, endpoint string, buf *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set(server.AuthHeader, a.encodeAuth())
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if a.config.GZip {
		req.Header.Set("Content-Encoding", "gzip")
	}
	return req, nil
}

// sendRequest send http request
func (a *agent) sendRequest(ctx context.Context, path string, method string, v any, rbody bool) {
	endpoint := fmt.Sprintf("http://%v/%s", a.config.Host, path)
	var buf bytes.Buffer
	a.logger.Debug("sendRequest endpoint: ", endpoint)
	if v != nil {
		a.write2Buffer(&buf, v)
	}
	req, err := a.getRequest(ctx, method, endpoint, &buf)
	if err != nil {
		a.logger.Errorf("agent.sendRequest req err: %s", err.Error())
		return
	}
	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Errorf("agent.sendRequest resp err: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	a.logger.Infof("sendRequest response code: %d", resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		a.logger.Errorf("sendRequest err: %s", err.Error())
	}
	if rbody {
		fmt.Printf("Response:\n %v", bytes.NewBuffer(bodyBytes).String())
	}

}

func (a *agent) register(ctx context.Context) {
	uri := "api/v1/user/register"
	m := &model.User{Login: a.config.Login, Password: a.config.Password}
	a.sendRequest(ctx, uri, http.MethodPost, m, false)
}

func (a *agent) health(ctx context.Context) {
	uri := "api/v1/private/health"
	m := &model.User{Login: a.config.Login, Password: a.config.Password}
	a.sendRequest(ctx, uri, http.MethodGet, m, false)
}

func (a *agent) listLogingWithPassword(ctx context.Context) {
	uri := "api/v1/private/loginwithpasswords"
	a.sendRequest(ctx, uri, http.MethodGet, nil, true)
}

func (a *agent) getLogingWithPassword(ctx context.Context) {
	uri := "api/v1/private/loginwithpasswords"
	a.sendRequest(ctx, uri, http.MethodGet, nil, true)
}

func (a *agent) addLogingWithPassword(ctx context.Context, m *model.LoginWithPassword) {
	uri := "api/v1/private/loginwithpassword"
	a.sendRequest(ctx, uri, http.MethodPut, m, true)
}

func (a *agent) deleteLogingWithPassword(ctx context.Context, id string) {
	uri := "api/v1/private/loginwithpassword/" + id
	a.sendRequest(ctx, uri, http.MethodDelete, nil, true)
}
func (a *agent) updateLogingWithPassword(ctx context.Context, m *model.LoginWithPassword) {
	uri := "api/v1/private/loginwithpassword"
	a.sendRequest(ctx, uri, http.MethodPost, m, true)
}

func (a *agent) list(ctx context.Context, target string) {
	switch target {
	case "l", "login", "password", "lp":
		a.listLogingWithPassword(ctx)
	case "c", "credit", "card", "cc", "creditcard":
		a.listLogingWithPassword(ctx)
	case "t", "text", "secrettext":
		a.listLogingWithPassword(ctx)
	case "f", "file", "secretfile":
		a.listLogingWithPassword(ctx)
	default:
		log.Fatalf("Unknown target: %s", target)
	}
}

func (a *agent) get(ctx context.Context, target string) {
	switch target {
	case "l", "login", "password", "lp":
		a.getLogingWithPassword(ctx)
	case "c", "credit", "card", "cc", "creditcard":
		a.listLogingWithPassword(ctx)
	case "t", "text", "secrettext":
		a.listLogingWithPassword(ctx)
	case "f", "file", "secretfile":
		a.listLogingWithPassword(ctx)
	default:
		log.Fatalf("Unknown target: %s", target)
	}
}

func (a *agent) add(ctx context.Context, target string) {
	switch target {
	case "l", "login", "password", "lp":
		login := input("Login")
		password := input("Password")
		m := &model.LoginWithPassword{Login: login, Password: password}
		a.addLogingWithPassword(ctx, m)
	case "c", "credit", "card", "cc", "creditcard":
		a.listLogingWithPassword(ctx)
	case "t", "text", "secrettext":
		a.listLogingWithPassword(ctx)
	case "f", "file", "secretfile":
		a.listLogingWithPassword(ctx)
	default:
		log.Fatalf("Unknown target: %s", target)
	}
}

func (a *agent) delete(ctx context.Context, target string) {
	switch target {
	case "l", "login", "password", "lp":
		a.listLogingWithPassword(ctx)
		id := input("Id of secret to delete")
		a.deleteLogingWithPassword(ctx, id)
	case "c", "credit", "card", "cc", "creditcard":
		a.listLogingWithPassword(ctx)
	case "t", "text", "secrettext":
		a.listLogingWithPassword(ctx)
	case "f", "file", "secretfile":
		a.listLogingWithPassword(ctx)
	default:
		log.Fatalf("Unknown target: %s", target)
	}
}

func (a *agent) update(ctx context.Context, target string) {
	switch target {
	case "l", "login", "password", "lp":
		a.listLogingWithPassword(ctx)
		id := input("Id of secret to update")
		login := input("Login")
		password := input("Password")
		meta := input("Meta")
		name := input("Name")
		m := &model.LoginWithPassword{Login: login, Password: password}
		m.ID, _ = strconv.Atoi(id)
		m.Name = name
		m.Password = password
		m.Meta = meta
		a.updateLogingWithPassword(ctx, m)
	case "c", "credit", "card", "cc", "creditcard":
		a.listLogingWithPassword(ctx)
	case "t", "text", "secrettext":
		a.listLogingWithPassword(ctx)
	case "f", "file", "secretfile":
		a.listLogingWithPassword(ctx)
	default:
		log.Fatalf("Unknown target: %s", target)
	}
}

func (a *agent) userInput() {
	ctx := context.Background()
	action := input("Action")
	a.logger.Infof("agent.userInput action: %s", action)
	if action == "register" {
		a.register(ctx)
		return
	}
	target := input("Type of secret")
	switch action {
	case "list":
		a.list(ctx, target)
	case "get":
		a.get(ctx, target)
	case "add":
		a.add(ctx, target)
	case "delete":
		a.delete(ctx, target)
	case "update":
		a.update(ctx, target)
	default:
		log.Fatalf("Unknown action: %s", action)
	}
	a.logger.Info("Done")
}
