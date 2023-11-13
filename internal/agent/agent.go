package agent

import (
	"bytes"
	"cenarius/internal/cache"
	"cenarius/internal/cache/filecache"
	"cenarius/internal/cache/mcache"
	"cenarius/internal/model"
	"cenarius/internal/server"
	"cenarius/internal/userinput"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
)

var errBadHTTPStatusCode = errors.New("bad http status code")

type agent struct {
	client     http.Client
	config     *Config
	logger     *logrus.Logger
	cache      cache.StoreCache
	store      cache.StoreCache
	onlineMode bool
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
	return a
}

// Start starts the agent
func (a *agent) Start() error {
	ctx := context.Background()
	a.logger.Info("Configuring logger")
	if err := a.configureLogger(); err != nil {
		return err
	}
	a.setKeyAndIV()
	a.logger.Info("Configuring store")
	if err := a.configureStore(); err != nil {
		return err
	}
	a.logger.Info("Reading cache from store")
	if err := a.readCache(); err != nil {
		return err
	}
	a.logger.Info("Checking the server availability")
	statusCode, err := a.ping(ctx)
	if err != nil {
		return err
	}
	if statusCode < 0 {
		a.onlineMode = false
	}
	if statusCode == http.StatusUnauthorized {
		a.register(ctx)
	}
	if err := a.updateCache(ctx); err != nil {
		return err
	}
	a.userInput()
	return nil
}

// Stop stops the agent
func (a *agent) Shutdown() {
	if err := a.saveCache(); err != nil {
		a.logger.Errorf("Unable to save cache: %s", err.Error())
	}
	a.store.Cache().Close()
	os.Exit(0)
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
func (a *agent) configureStore() error {
	f, err := os.OpenFile(a.config.CacheFile, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		a.logger.Errorf("agent.configureStore Failed to open file %v", a.config.CacheFile)
		return err
	}
	a.store = filecache.New(f)
	a.cache = mcache.New()
	return nil
}

func (a *agent) setKeyAndIV() {
	a.logger.Debug("agent.getKeyAndIV is working")
	key := a.config.Login + a.config.Password
	a.logger.Debug("agent.getKeyAndIV key: ", key)
	if len(key) < 32 {
		for i := 1; len(key) < 32; i++ {
			key += strconv.Itoa(i)
		}
	}
	key = key[0:32]
	iv := key[0:16]
	a.logger.Debug("agent.getKeyAndIV Key, vector: ", key, iv)
	a.config.SecretKey = key
	a.config.SecretIV = iv
}
func (a *agent) readCache() error {
	c, err := a.store.Cache().Get()
	if err != nil {
		return err
	}
	a.logger.Debug("agent.readCache cache: ", c)
	if c != nil {
		if err := a.cache.Cache().Save(c); err != nil {
			return err
		}
	}
	return nil
}

func (a *agent) getSecretsWrapper(ctx context.Context, uri string, v any) error {
	data, s, err := a.sendRequest2(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return err
	}
	if s != http.StatusOK && s != http.StatusCreated {
		a.logger.Errorf("Faliled to get %s: %d", uri, s)
		return errBadHTTPStatusCode
	}
	a.logger.Debugf("Got from %s: %v", uri, string(data))
	if err := json.Unmarshal(data, &v); err != nil {
		a.logger.Errorf("agent.getSecrets unmarshal json failed %v: %v", string(data), err)
		return err
	}
	return nil
}

func (a *agent) getSecrets(ctx context.Context) (*model.SecretCache, error) {
	var cache = &model.SecretCache{}
	if err := a.getSecretsWrapper(ctx, logingWithPasswordGetURI, &cache.LoginWithPasswords); err != nil {
		return nil, err
	}
	if err := a.getSecretsWrapper(ctx, creditCardGetURI, &cache.CreditCards); err != nil {
		return nil, err
	}
	if err := a.getSecretsWrapper(ctx, secretTextGetURI, &cache.SecretTexts); err != nil {
		return nil, err
	}
	if err := a.getSecretsWrapper(ctx, secretFileGetURI, &cache.SecretFiles); err != nil {
		return nil, err
	}
	a.logger.Debugf("Got new cache from server: %v", cache)
	return cache, nil
}

func (a *agent) updateCache(ctx context.Context) error {
	cache, err := a.getSecrets(ctx)
	if err != nil {
		return err
	}
	a.logger.Debugf("Got secret: %v", cache)
	if err := cache.Encrypt(a.config.SecretKey, a.config.SecretIV); err != nil {
		return err
	}
	if err := a.store.Cache().Save(cache); err != nil {
		return err
	}
	return nil
}

func (a *agent) saveCache() error {
	c, err := a.cache.Cache().Get()
	if err != nil {
		return err
	}
	if err := c.Encrypt(a.config.SecretKey, a.config.SecretIV); err != nil {
		return err
	}
	if err := a.store.Cache().Save(c); err != nil {
		return err
	}
	return nil
}

func (a *agent) printEncryptedSecret(i model.Encrypter) error {
	if err := i.Decrypt(a.config.SecretKey, a.config.SecretIV); err != nil {
		a.logger.Errorf("agent.printEncryptedSecret failed to decrypt %v : %v", i, err.Error())
		return err
	}
	fmt.Println(i)
	if err := i.Encrypt(a.config.SecretKey, a.config.SecretIV); err != nil {
		a.logger.Errorf("agent.printEncryptedSecret failed to encrypt %v : %v", i, err.Error())
		return err
	}
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
			a.logger.Errorf("agent.write2Buffer err: %s", err.Error())
			return
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
		a.logger.Errorf("agent.sendRequest req err: %s", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Errorf("agent.sendRequest resp err: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	a.logger.Debugf("sendRequest response code: %d", resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		a.logger.Errorf("sendRequest err: %s", err.Error())
		return
	}
	if rbody {
		fmt.Printf("Response:\n %v", bytes.NewBuffer(bodyBytes).String())
	}
}

// sendRequest send http request
func (a *agent) sendRequest2(ctx context.Context, path string, method string, v any) ([]byte, int, error) {
	endpoint := a.geHTTPtURL(path)
	var buf bytes.Buffer
	a.logger.Debug("sendRequest endpoint: ", endpoint)
	if v != nil {
		a.write2Buffer(&buf, v)
	}
	req, err := a.getRequest(ctx, method, endpoint, &buf)
	if err != nil {
		a.logger.Errorf("agent.sendRequest req err: %s", err.Error())
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Errorf("agent.sendRequest resp err: %s", err.Error())
		return nil, 0, err
	}
	defer resp.Body.Close()
	a.logger.Debugf("sendRequest response code: %d", resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		a.logger.Errorf("sendRequest err: %s", err.Error())
		return nil, 0, err
	}
	return bodyBytes, resp.StatusCode, nil
}

func (a *agent) register(ctx context.Context) {
	m := &model.User{Login: a.config.Login, Password: a.config.Password}
	a.sendRequest(ctx, registerURI, http.MethodPost, m, false)
}

func (a *agent) ping(ctx context.Context) (int, error) {
	m := &model.User{Login: a.config.Login, Password: a.config.Password}
	_, s, err := a.sendRequest2(ctx, pingURI, http.MethodGet, m)
	if err != nil {
		a.logger.Errorf("agent.ping error: %s", err.Error())
		return 0, err
	}
	return s, nil
}

func (a *agent) listLogingWithPassword(ctx context.Context) {
	fmt.Println("You Logins and Passwords: ")
	cache, err := a.cache.Cache().Get()
	if err != nil {
		a.logger.Error(err.Error())
		return
	}
	l := make([]*model.LoginWithPassword, len(cache.LoginWithPasswords))
	copy(l, cache.LoginWithPasswords)
	for _, i := range l {
		fmt.Println(i)
	}
}

func (a *agent) getLogingWithPassword(ctx context.Context, id int) {
	fmt.Printf("You Login and Password with id: %d\n", id)
	cache, err := a.cache.Cache().Get()
	if err != nil {
		a.logger.Error(err.Error())
		return
	}
	for _, i := range cache.LoginWithPasswords {
		if i.ID == id {
			if err := a.printEncryptedSecret(i); err != nil {
				return
			}
		}
	}
}

func (a *agent) addLogingWithPassword(ctx context.Context, m *model.LoginWithPassword) {
	a.sendRequest(ctx, logingWithPasswordBodyURI, http.MethodPost, m, true)
}

func (a *agent) updateLogingWithPassword(ctx context.Context, m *model.LoginWithPassword) {
	a.sendRequest(ctx, logingWithPasswordBodyURI, http.MethodPut, m, true)
}

func (a *agent) deleteLogingWithPassword(ctx context.Context, id int) {
	uri := fmt.Sprintf("%s/%s", logingWithPasswordBodyURI, strconv.Itoa(id))
	a.sendRequest(ctx, uri, http.MethodDelete, nil, true)
}

func (a *agent) listCreditCard(ctx context.Context) {
	fmt.Println("You Credit Cards: ")
	cache, err := a.cache.Cache().Get()
	if err != nil {
		a.logger.Error(err.Error())
		return
	}
	for _, i := range cache.CreditCards {
		fmt.Println(i)
	}
}

func (a *agent) getCreditCard(ctx context.Context, id int) {
	fmt.Printf("You Credit Cards with id: %d\n", id)
	cache, err := a.cache.Cache().Get()
	if err != nil {
		a.logger.Error(err.Error())
		return
	}
	for _, i := range cache.CreditCards {
		if i.ID == id {
			if err := a.printEncryptedSecret(i); err != nil {
				return
			}
		}
	}
}

func (a *agent) addCreditCard(ctx context.Context, m *model.CreditCard) {
	a.sendRequest(ctx, creditCardBodyURI, http.MethodPost, m, true)
}

func (a *agent) updateCreditCard(ctx context.Context, m *model.CreditCard) {
	a.sendRequest(ctx, creditCardBodyURI, http.MethodPut, m, true)
}

func (a *agent) deleteCreditCard(ctx context.Context, id int) {
	uri := fmt.Sprintf("%s/%s", creditCardBodyURI, strconv.Itoa(id))
	a.sendRequest(ctx, uri, http.MethodDelete, nil, true)
}

func (a *agent) listSecretText(ctx context.Context) {
	fmt.Println("You Secret Texts: ")
	cache, err := a.cache.Cache().Get()
	if err != nil {
		a.logger.Error(err.Error())
		return
	}
	for _, i := range cache.SecretTexts {
		if err := i.Encrypt(a.config.SecretKey, a.config.SecretIV); err != nil {
			return
		}
		fmt.Println(i)
		if err := i.Decrypt(a.config.SecretKey, a.config.SecretIV); err != nil {
			return
		}
	}
}

func (a *agent) getSecretText(ctx context.Context, id int) {
	fmt.Printf("You Secret Texts with id: %d\n", id)
	cache, err := a.cache.Cache().Get()
	if err != nil {
		a.logger.Error(err.Error())
		return
	}
	for _, i := range cache.SecretTexts {
		if i.ID == id {
			if err := a.printEncryptedSecret(i); err != nil {
				return
			}
		}
	}
}

func (a *agent) addSecretText(ctx context.Context, m *model.SecretText) {
	a.sendRequest(ctx, secretTextBodyURI, http.MethodPost, m, true)
}

func (a *agent) deleteSecretText(ctx context.Context, id int) {
	uri := fmt.Sprintf("%s/%s", secretTextBodyURI, strconv.Itoa(id))
	a.sendRequest(ctx, uri, http.MethodDelete, nil, true)
}
func (a *agent) updateSecretText(ctx context.Context, m *model.SecretText) {
	a.sendRequest(ctx, secretTextBodyURI, http.MethodPut, m, true)
}

func (a *agent) listSecretFile(ctx context.Context) {
	fmt.Println("You Secret Files: ")
	cache, err := a.cache.Cache().Get()
	if err != nil {
		a.logger.Error(err.Error())
		return
	}
	for _, i := range cache.SecretFiles {
		i.Encrypt(a.config.SecretKey, a.config.SecretIV)
		fmt.Println(i)
		i.Decrypt(a.config.SecretKey, a.config.SecretIV)
	}
}

func (a *agent) getSecretFile(ctx context.Context, id string) {
	uri := fmt.Sprintf("%s/%s", secretFileBodyURI, id)
	endpoint := a.geHTTPtURL(uri)
	req, err := a.getRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		a.logger.Errorf("agent.getSecretFile req err: %s", err.Error())
		return
	}
	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Errorf("agent.getSecretFile resp err: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		a.logger.Errorf("bad status: %s", resp.Status)
		return
	}
	out, err := os.Create("SecretFile_" + id)
	if err != nil {
		a.logger.Errorf("Can't create local file: %s", err.Error())
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		a.logger.Errorf("Can't copy reposnse to local file: %s", err.Error())
		return
	}
}

func (a *agent) uploadSecretFile(ctx context.Context, m *model.SecretFile) {
	file, err := os.Open(m.Path)
	if err != nil {
		a.logger.Errorf("The file doesn't exist: %s", m.Path)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("secretFile", filepath.Base(file.Name()))
	if err != nil {
		a.logger.Errorf("agent.addSecretFile CreateFormFile %s: %s", file.Name(), err.Error())
		return
	}
	if _, err := io.Copy(part, file); err != nil {
		a.logger.Errorf("Failed copy part of file: %s", err.Error())
		return
	}
	writer.Close()
	endpoint := a.geHTTPtURL(secretFileBodyURI)
	req, err := a.getRequest(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		a.logger.Errorf("agent.addSecretFile req err: %s", err.Error())
		return
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Errorf("agent.addSecretFile resp err: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	responseM := &model.SecretFile{}
	if err := json.NewDecoder(resp.Body).Decode(responseM); err != nil {
		a.logger.Errorf("Unable to parse resp body in uploadSecretFile: %v", err)
		return
	}
	a.logger.Infof("agent.uploadSecretFile uploaded: %v", responseM)
	m.ID = responseM.ID
	m.UserID = responseM.ID
	a.updateSecretFile(ctx, m)
}

func (a *agent) updateSecretFile(ctx context.Context, m *model.SecretFile) {
	a.sendRequest(ctx, secretFileBodyURI, http.MethodPut, m, true)
}

func (a *agent) deleteSecretFile(ctx context.Context, id int) {
	uri := fmt.Sprintf("%s/%s", secretFileBodyURI, strconv.Itoa(id))
	a.sendRequest(ctx, uri, http.MethodDelete, nil, true)
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
		a.logger.Errorf("Unknown target: %s", target)
	}
}

func (a *agent) get(ctx context.Context, target string) {
	switch target {
	case "l", "login", "password", "lp":
		a.listLogingWithPassword(ctx)
		id := userinput.InputID()
		a.getLogingWithPassword(ctx, id)
	case "c", "credit", "card", "cc", "creditcard":
		a.listCreditCard(ctx)
		id := userinput.InputID()
		a.getCreditCard(ctx, id)
	case "t", "text", "secrettext":
		a.listSecretText(ctx)
		id := userinput.InputID()
		a.getSecretText(ctx, id)
	case "f", "file", "secretfile":
		a.listSecretFile(ctx)
		id := userinput.InputID()
		a.getSecretFile(ctx, strconv.Itoa(id))
	default:
		a.logger.Errorf("Unknown target: %s", target)
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
		a.logger.Errorf("Unknown target: %s", target)
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
		a.logger.Errorf("Unknown target: %s", target)
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
		a.logger.Errorf("Unknown target: %s", target)
	}
}

func (a *agent) userInput() {
	ctx := context.Background()
	action := userinput.Input("Action: (r|register) (l|list) (g|get) (a|add) (d|delete) (u|update)")
	a.logger.Infof("agent.userInput action: %s", action)
	if action == "register" || action == "r" {
		a.register(ctx)
		return
	}
	target := userinput.Input("Type of secret you want to operate: (l|login|password|lp) (c|credit|card|cc|creditcard) (t|text|secrettext) (f|file|secretfile)")
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
		a.logger.Errorf("Unknown action: %s", action)
	}
	if err := a.updateCache(ctx); err != nil {
		a.logger.Error(err)
	}
	a.logger.Info("Done")
}
