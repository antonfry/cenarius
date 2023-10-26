package agent

import (
	"bytes"
	"cenarius/internal/agent/userinput"
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

func (a *agent) deleteLogingWithPassword(ctx context.Context, id int) {
	uri := "api/v1/private/loginwithpassword/" + strconv.Itoa(id)
	a.sendRequest(ctx, uri, http.MethodDelete, nil, true)
}
func (a *agent) updateLogingWithPassword(ctx context.Context, m *model.LoginWithPassword) {
	uri := "api/v1/private/loginwithpassword"
	a.sendRequest(ctx, uri, http.MethodPost, m, true)
}

func (a *agent) listCreditCard(ctx context.Context) {
	uri := "api/v1/private/creditcards"
	a.sendRequest(ctx, uri, http.MethodGet, nil, true)
}

func (a *agent) getCreditCard(ctx context.Context) {
	uri := "api/v1/private/creditcard"
	a.sendRequest(ctx, uri, http.MethodGet, nil, true)
}

func (a *agent) addCreditCard(ctx context.Context, m *model.CreditCard) {
	uri := "api/v1/private/creditcard"
	a.sendRequest(ctx, uri, http.MethodPut, m, true)
}

func (a *agent) deleteCreditCard(ctx context.Context, id int) {
	uri := "api/v1/private/creditcard/" + strconv.Itoa(id)
	a.sendRequest(ctx, uri, http.MethodDelete, nil, true)
}
func (a *agent) updateCreditCard(ctx context.Context, m *model.CreditCard) {
	uri := "api/v1/private/creditcard"
	a.sendRequest(ctx, uri, http.MethodPost, m, true)
}

func (a *agent) listSecretText(ctx context.Context) {
	uri := "api/v1/private/secrettexts"
	a.sendRequest(ctx, uri, http.MethodGet, nil, true)
}

func (a *agent) getSecretText(ctx context.Context) {
	uri := "api/v1/private/secrettext"
	a.sendRequest(ctx, uri, http.MethodGet, nil, true)
}

func (a *agent) addSecretText(ctx context.Context, m *model.SecretText) {
	uri := "api/v1/private/secrettext"
	a.sendRequest(ctx, uri, http.MethodPut, m, true)
}

func (a *agent) deleteSecretText(ctx context.Context, id int) {
	uri := "api/v1/private/secrettext/" + strconv.Itoa(id)
	a.sendRequest(ctx, uri, http.MethodDelete, nil, true)
}
func (a *agent) updateSecretText(ctx context.Context, m *model.SecretText) {
	uri := "api/v1/private/secrettext"
	a.sendRequest(ctx, uri, http.MethodPost, m, true)
}

func (a *agent) listSecretFile(ctx context.Context) {
	uri := "api/v1/private/secretfiles"
	a.sendRequest(ctx, uri, http.MethodGet, nil, true)
}

func (a *agent) getSecretFile(ctx context.Context) {
	uri := "api/v1/private/secretfile"
	a.sendRequest(ctx, uri, http.MethodGet, nil, true)
}

func (a *agent) addSecretFile(ctx context.Context, m *model.SecretFile) {
	uri := "api/v1/private/secretfile"
	a.sendRequest(ctx, uri, http.MethodPut, m, true)
}

func (a *agent) deleteSecretFile(ctx context.Context, id int) {
	uri := "api/v1/private/secretfile/" + strconv.Itoa(id)
	a.sendRequest(ctx, uri, http.MethodDelete, nil, true)
}
func (a *agent) updateSecretFile(ctx context.Context, m *model.SecretFile) {
	uri := "api/v1/private/secretfile"
	a.sendRequest(ctx, uri, http.MethodPost, m, true)
}

func (a *agent) list(ctx context.Context, target string) {
	switch target {
	case "l", "login", "password", "lp":
		a.listLogingWithPassword(ctx)
	case "c", "credit", "card", "cc", "creditcard":
		a.listLogingWithPassword(ctx)
	case "t", "text", "secrettext":
		a.listSecretText(ctx)
	case "f", "file", "secretfile":
		a.listSecretFile(ctx)
	default:
		log.Fatalf("Unknown target: %s", target)
	}
}

func (a *agent) get(ctx context.Context, target string) {
	switch target {
	case "l", "login", "password", "lp":
		a.getLogingWithPassword(ctx)
	case "c", "credit", "card", "cc", "creditcard":
		a.getCreditCard(ctx)
	case "t", "text", "secrettext":
		a.getSecretText(ctx)
	case "f", "file", "secretfile":
		a.getSecretFile(ctx)
	default:
		log.Fatalf("Unknown target: %s", target)
	}
}

func (a *agent) add(ctx context.Context, target string) {
	switch target {
	case "l", "login", "password", "lp":
		m := userinput.InputLoginWithPassword()
		a.addLogingWithPassword(ctx, m)
	case "c", "credit", "card", "cc", "creditcard":
		m := userinput.InputCreditCard()
		a.addCreditCard(ctx, m)
	case "t", "text", "secrettext":
		m := userinput.InputSecretText()
		a.addSecretText(ctx, m)
	case "f", "file", "secretfile":
		m := userinput.InputSecretFile()
		a.addSecretFile(ctx, m)
	default:
		log.Fatalf("Unknown target: %s", target)
	}
}

func (a *agent) delete(ctx context.Context, target string) {
	switch target {
	case "l", "login", "password", "lp":
		a.listLogingWithPassword(ctx)
		id := userinput.InputId()
		a.deleteLogingWithPassword(ctx, id)
	case "c", "credit", "card", "cc", "creditcard":
		a.listCreditCard(ctx)
		id := userinput.InputId()
		a.deleteCreditCard(ctx, id)
	case "t", "text", "secrettext":
		a.listSecretText(ctx)
		id := userinput.InputId()
		a.deleteSecretText(ctx, id)
	case "f", "file", "secretfile":
		a.listSecretFile(ctx)
		id := userinput.InputId()
		a.deleteSecretFile(ctx, id)
	default:
		log.Fatalf("Unknown target: %s", target)
	}
}

func (a *agent) update(ctx context.Context, target string) {
	switch target {
	case "l", "login", "password", "lp":
		a.listLogingWithPassword(ctx)
		id := userinput.InputId()
		m := userinput.InputLoginWithPassword()
		m.ID = id
		a.updateLogingWithPassword(ctx, m)
	case "c", "credit", "card", "cc", "creditcard":
		a.listCreditCard(ctx)
		id := userinput.InputId()
		m := userinput.InputCreditCard()
		m.ID = id
		a.updateCreditCard(ctx, m)
	case "t", "text", "secrettext":
		a.listSecretText(ctx)
		id := userinput.InputId()
		m := userinput.InputSecretText()
		m.ID = id
		a.updateSecretText(ctx, m)
	case "f", "file", "secretfile":
		a.listSecretFile(ctx)
		id := userinput.InputId()
		m := userinput.InputSecretFile()
		m.ID = id
		a.updateSecretFile(ctx, m)
	default:
		log.Fatalf("Unknown target: %s", target)
	}
}

func (a *agent) userInput() {
	ctx := context.Background()
	action := userinput.Input("Action")
	a.logger.Infof("agent.userInput action: %s", action)
	if action == "register" {
		a.register(ctx)
		return
	}
	target := userinput.Input("Type of secret")
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
