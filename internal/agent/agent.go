package agent

import (
	"bytes"
	"cenarius/internal/cache"
	"cenarius/internal/cache/filecache"
	"cenarius/internal/model"
	"cenarius/internal/server"
	"cenarius/internal/userinput"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
)

type agent struct {
	client http.Client
	config *Config
	logger *logrus.Logger
	store  cache.StoreCache
}

// NewServer returns new server object
func NewAgent(config *Config) *agent {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	a := &agent{
		client: http.Client{Transport: tr},
		config: config,
		logger: logrus.New(),
	}
	f, err := os.OpenFile(a.config.CacheFile, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		a.logger.Errorf("openFile: Failed to open file %v", a.config.CacheFile)
		return nil
	}
	a.store = filecache.New(f)
	return a
}

// Start starts the agent
func (a *agent) Start() error {
	ctx := context.Background()
	if err := a.configureLogger(); err != nil {
		return err
	}
	a.ping(ctx)
	a.userInput()
	return nil
}

// Stop stops the agent
func (a *agent) Shutdown() {
	a.store.Cache().Close()
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
	a.logger.Debugf("write2Buffer jsonData: %s", string(jsonData))
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
		if _, err := io.WriteString(buf, string(jsonData)); err != nil {
			a.logger.Fatalf("agent.write2Buffer err: %s", err.Error())
		}
	}
}

func (a *agent) encodeAuth() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s %s", a.config.Login, a.config.Password)))
}

func (a *agent) getRequest(ctx context.Context, method string, endpoint string, buf *bytes.Buffer) (*http.Request, error) {
	var req *http.Request
	var err error
	if buf == nil {
		req, err = http.NewRequestWithContext(ctx, method, endpoint, nil)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, endpoint, buf)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set(server.AuthHeader, a.encodeAuth())
	if a.config.GZip {
		req.Header.Set("Content-Encoding", "gzip")
	}
	return req, nil
}

func (a *agent) geHTTPtURL(path string) string {
	return fmt.Sprintf("https://%v/%s", a.config.Host, path)
}

// sendRequest send http request
func (a *agent) sendRequest(ctx context.Context, path string, method string, v any, rbody bool) {
	endpoint := a.geHTTPtURL(path)
	var buf bytes.Buffer
	a.logger.Debug("sendRequest endpoint: ", endpoint)
	if v != nil {
		a.write2Buffer(&buf, v)
	}
	req, err := a.getRequest(ctx, method, endpoint, &buf)
	if err != nil {
		a.logger.Fatalf("agent.sendRequest req err: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Fatalf("agent.sendRequest resp err: %s", err.Error())
	}
	defer resp.Body.Close()
	a.logger.Debugf("sendRequest response code: %d", resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		a.logger.Fatalf("sendRequest err: %s", err.Error())
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

func (a *agent) ping(ctx context.Context) {
	uri := "api/v1/private/ping"
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
	a.sendRequest(ctx, uri, http.MethodPost, m, true)
}

func (a *agent) deleteLogingWithPassword(ctx context.Context, id int) {
	uri := "api/v1/private/loginwithpassword/" + strconv.Itoa(id)
	a.sendRequest(ctx, uri, http.MethodDelete, nil, true)
}
func (a *agent) updateLogingWithPassword(ctx context.Context, m *model.LoginWithPassword) {
	uri := "api/v1/private/loginwithpassword"
	a.sendRequest(ctx, uri, http.MethodPut, m, true)
}

func (a *agent) listCreditCard(ctx context.Context) {
	uri := "api/v1/private/creditcards"
	a.sendRequest(ctx, uri, http.MethodGet, nil, true)
}

func (a *agent) getCreditCard(ctx context.Context) {
	uri := "api/v1/private/creditcards"
	a.sendRequest(ctx, uri, http.MethodGet, nil, true)
}

func (a *agent) addCreditCard(ctx context.Context, m *model.CreditCard) {
	uri := "api/v1/private/creditcard"
	a.sendRequest(ctx, uri, http.MethodPost, m, true)
}

func (a *agent) deleteCreditCard(ctx context.Context, id int) {
	uri := "api/v1/private/creditcard/" + strconv.Itoa(id)
	a.sendRequest(ctx, uri, http.MethodDelete, nil, true)
}
func (a *agent) updateCreditCard(ctx context.Context, m *model.CreditCard) {
	uri := "api/v1/private/creditcard"
	a.sendRequest(ctx, uri, http.MethodPut, m, true)
}

func (a *agent) listSecretText(ctx context.Context) {
	uri := "api/v1/private/secrettexts"
	a.sendRequest(ctx, uri, http.MethodGet, nil, true)
}

func (a *agent) getSecretText(ctx context.Context) {
	uri := "api/v1/private/secrettexts"
	a.sendRequest(ctx, uri, http.MethodGet, nil, true)
}

func (a *agent) addSecretText(ctx context.Context, m *model.SecretText) {
	uri := "api/v1/private/secrettext"
	a.sendRequest(ctx, uri, http.MethodPost, m, true)
}

func (a *agent) deleteSecretText(ctx context.Context, id int) {
	uri := "api/v1/private/secrettext/" + strconv.Itoa(id)
	a.sendRequest(ctx, uri, http.MethodDelete, nil, true)
}
func (a *agent) updateSecretText(ctx context.Context, m *model.SecretText) {
	uri := "api/v1/private/secrettext"
	a.sendRequest(ctx, uri, http.MethodPut, m, true)
}

func (a *agent) listSecretFile(ctx context.Context) {
	uri := "api/v1/private/secretfiles"
	a.sendRequest(ctx, uri, http.MethodGet, nil, true)
}

func (a *agent) getSecretFile(ctx context.Context, id string) {
	uri := "api/v1/private/secretfile/" + id
	endpoint := a.geHTTPtURL(uri)

	// req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	// if err != nil {
	// 	a.logger.Fatalf("agent.getSecretFile create request err: %s", err.Error())
	// }
	// req.Header.Set(server.AuthHeader, a.encodeAuth())

	req, err := a.getRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		a.logger.Fatalf("agent.getSecretFile req err: %s", err.Error())
	}
	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Fatalf("agent.getSecretFile resp err: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		a.logger.Fatalf("bad status: %s", resp.Status)
	}
	out, err := os.Create("SecretFile_" + id)
	if err != nil {
		a.logger.Fatalf("Can't create local file: %s", err.Error())
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		a.logger.Fatalf("Can't copy reposnse to local file: %s", err.Error())
	}
}

func (a *agent) uploadSecretFile(ctx context.Context, m *model.SecretFile) {
	uri := "api/v1/private/secretfile"
	file, err := os.Open(m.Path)
	if err != nil {
		a.logger.Fatalf("The file doesn't exist: %s", m.Path)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("secretFile", filepath.Base(file.Name()))
	if err != nil {
		a.logger.Fatalf("agent.addSecretFile CreateFormFile %s: %s", file.Name(), err.Error())
	}
	if _, err := io.Copy(part, file); err != nil {
		a.logger.Fatalf("Failed copy part of file: %s", err.Error())
	}
	writer.Close()
	endpoint := a.geHTTPtURL(uri)
	req, err := a.getRequest(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		a.logger.Fatalf("agent.addSecretFile req err: %s", err.Error())
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Fatalf("agent.addSecretFile resp err: %s", err.Error())
	}
	defer resp.Body.Close()
	responseM := &model.SecretFile{}
	if err := json.NewDecoder(resp.Body).Decode(responseM); err != nil {
		a.logger.Fatalf("Unable to parse resp body in uploadSecretFile: %v", err)
	}
	a.logger.Infof("agent.uploadSecretFile uploaded: %v", responseM)
	m.ID = responseM.ID
	m.UserID = responseM.ID
	a.updateSecretFile(ctx, m)
}

func (a *agent) deleteSecretFile(ctx context.Context, id int) {
	uri := "api/v1/private/secretfile/" + strconv.Itoa(id)
	a.sendRequest(ctx, uri, http.MethodDelete, nil, true)
}
func (a *agent) updateSecretFile(ctx context.Context, m *model.SecretFile) {
	uri := "api/v1/private/secretfile"
	a.sendRequest(ctx, uri, http.MethodPut, m, true)
}

func (a *agent) list(ctx context.Context, target string) {
	switch target {
	case "l", "login", "password", "lp":
		a.listLogingWithPassword(ctx)
	case "c", "credit", "card", "cc", "creditcard":
		a.listCreditCard(ctx)
	case "t", "text", "secrettext":
		a.listSecretText(ctx)
	case "f", "file", "secretfile":
		a.listSecretFile(ctx)
	default:
		a.logger.Fatalf("Unknown target: %s", target)
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
		a.listSecretFile(ctx)
		id := userinput.InputID()
		a.getSecretFile(ctx, strconv.Itoa(id))
	default:
		a.logger.Fatalf("Unknown target: %s", target)
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
		m := userinput.InputSecretFile(true)
		a.uploadSecretFile(ctx, m)
	default:
		a.logger.Fatalf("Unknown target: %s", target)
	}
}

func (a *agent) delete(ctx context.Context, target string) {
	switch target {
	case "l", "login", "password", "lp":
		a.listLogingWithPassword(ctx)
		id := userinput.InputID()
		a.deleteLogingWithPassword(ctx, id)
	case "c", "credit", "card", "cc", "creditcard":
		a.listCreditCard(ctx)
		id := userinput.InputID()
		a.deleteCreditCard(ctx, id)
	case "t", "text", "secrettext":
		a.listSecretText(ctx)
		id := userinput.InputID()
		a.deleteSecretText(ctx, id)
	case "f", "file", "secretfile":
		a.listSecretFile(ctx)
		id := userinput.InputID()
		a.deleteSecretFile(ctx, id)
	default:
		a.logger.Fatalf("Unknown target: %s", target)
	}
}

func (a *agent) update(ctx context.Context, target string) {
	switch target {
	case "l", "login", "password", "lp":
		a.listLogingWithPassword(ctx)
		id := userinput.InputID()
		m := userinput.InputLoginWithPassword()
		m.ID = id
		a.updateLogingWithPassword(ctx, m)
	case "c", "credit", "card", "cc", "creditcard":
		a.listCreditCard(ctx)
		id := userinput.InputID()
		m := userinput.InputCreditCard()
		m.ID = id
		a.updateCreditCard(ctx, m)
	case "t", "text", "secrettext":
		a.listSecretText(ctx)
		id := userinput.InputID()
		m := userinput.InputSecretText()
		m.ID = id
		a.updateSecretText(ctx, m)
	case "f", "file", "secretfile":
		a.listSecretFile(ctx)
		id := userinput.InputID()
		m := userinput.InputSecretFile(false)
		m.ID = id
		a.updateSecretFile(ctx, m)
	default:
		a.logger.Fatalf("Unknown target: %s", target)
	}
}

func (a *agent) userInput() {
	ctx := context.Background()
	action := userinput.Input("Action")
	a.logger.Infof("agent.userInput action: %s", action)
	if action == "register" || action == "r" {
		a.register(ctx)
		return
	}
	target := userinput.Input("Type of secret")
	switch action {
	case "list", "l":
		a.list(ctx, target)
	case "get", "g":
		a.get(ctx, target)
	case "add", "a":
		a.add(ctx, target)
	case "delete", "d":
		a.delete(ctx, target)
	case "update", "u":
		a.update(ctx, target)
	default:
		a.logger.Fatalf("Unknown action: %s", action)
	}
	a.logger.Info("Done")
}
