package server

import (
	"cenarius/internal/model"
	"cenarius/internal/store"
	"cenarius/internal/store/sqlstore"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang-migrate/migrate"
	log "github.com/sirupsen/logrus"
)

type ctxKey int8

const (
	AuthHeader        = "X-Cenarius-Token"
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

var ErrUnableToGetUserFromRequest = errors.New("unable to get user from request context")

// server server main struct
type server struct {
	config     *Config
	logger     *log.Logger
	HTTPServer *http.Server
	router     *chi.Mux
	store      store.Store
}

// NewServer returns new server object
func NewServer(config *Config) *server {
	s := &server{
		config:     config,
		logger:     log.New(),
		HTTPServer: &http.Server{Addr: config.Bind},
	}
	if err := s.configureLogger(); err != nil {
		log.Fatalf("Can't configure logger: %s", err.Error())
	}
	if err := s.configureStore(); err != nil {
		log.Fatalf("Can't configure store: %s", err.Error())
	}

	return s
}

// StartHTTPServer starts GRPC Server
func (s *server) StartHTTPServer() error {
	s.logger.Infof("Starting HTTP server with config: %v\n", s.config)
	s.router = chi.NewRouter()
	s.configureRouter()
	if err := s.HTTPServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	s.logger.Infof("HTTP server stopped with config: %v\n", s.config)
	return nil
}

func (s *server) StopHTTPServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.HTTPServer.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		s.logger.Errorf("HTTP server Shutdown Error: %v", err)
	}
}

// Start starts the server
func (s *server) Start() error {
	err := s.StartHTTPServer()
	if err != nil {
		return err
	}
	return nil
}

// Shutdown shutdowns the server
func (s *server) Shutdown() {
	s.logger.Info("Shuting down...")
	s.StopHTTPServer()
	s.store.Close()
	s.logger.Info("Done ShutDown")
}

// configureLogger configures logger
func (s *server) configureLogger() error {
	level, err := log.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

// configureStore configures store
func (s *server) configureStore() error {
	var conn *sql.DB
	var err error
	for i := 1; i <= 10; i++ {
		conn, err = sqlstore.NewPGConn(s.config.DatabaseDsn)
		if i > 10 {
			return err
		}
		if err != nil && i <= 10 {
			s.logger.Errorf("Unable to connect to the database with: %v", s.config.DatabaseDsn)
			s.logger.Error(err)
			time.Sleep(time.Second * 2)
			continue
		} else {
			break
		}
	}
	if err := sqlstore.MigrateSQL(conn, s.config.MigrationPath); err != nil && err.Error() != migrate.ErrNoChange.Error() {
		s.logger.Error("Migration fail: ", err.Error())
		return err
	}
	s.store = sqlstore.NewStore(conn)
	return nil
}

func (s *server) userRegister(ctx context.Context, u *model.User) (*model.User, int, error) {
	if _, err := s.store.User().FindByLogin(ctx, u.Login); err == nil {
		s.logger.Errorf("User already exist")
		return nil, http.StatusConflict, store.ErrNotAuthenticated
	}
	if err := s.store.User().Create(ctx, u); err != nil {
		s.logger.Errorf("Failed to create user %v: %v", u, err)
		return nil, http.StatusBadRequest, err
	}
	s.logger.Debugf("User created: %v", u)
	return u, http.StatusAccepted, nil
}

func (s *server) userLogin(ctx context.Context, u *model.User) (*model.User, error) {
	storageUser, err := s.store.User().FindByLogin(ctx, u.Login)
	if err != nil {
		s.logger.Errorf("Unknown login: %s", u.Login)
		return nil, err
	}
	if !storageUser.ComparePassword(u.Password) {
		s.logger.Errorf("Incorrect Password: %v", u.Password)
		return nil, store.ErrIncorrectPassword
	}
	u.ID = storageUser.ID
	u.EncryptedPassword = storageUser.EncryptedPassword
	u.Sanitaze()
	return u, nil
}

func (s *server) addLoginWithPassword(ctx context.Context, m *model.LoginWithPassword, key, iv string) (*model.LoginWithPassword, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if err := m.Encrypt(key, iv); err != nil {
		return nil, err
	}
	if err := s.store.LoginWithPassword().Add(ctx, m); err != nil {
		s.logger.Errorf("Failed to add LoginWithPassword %v: %v", m, err)
		return nil, err
	}
	s.logger.Debugf("LoginWithPassword created: %v", m)
	return m, nil
}

func (s *server) addCreditCard(ctx context.Context, m *model.CreditCard, key, iv string) (*model.CreditCard, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if err := m.Encrypt(key, iv); err != nil {
		return nil, err
	}
	if err := s.store.CreditCard().Add(ctx, m); err != nil {
		s.logger.Errorf("Failed to add CreditCard %v: %v", m, err)
		return nil, err
	}
	s.logger.Debugf("CreditCard created: %v", m)
	return m, nil
}

func (s *server) addSecretText(ctx context.Context, m *model.SecretText, key, iv string) (*model.SecretText, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if err := m.Encrypt(key, iv); err != nil {
		return nil, err
	}
	if err := s.store.SecretText().Add(ctx, m); err != nil {
		s.logger.Errorf("Failed to add SecretText %v: %v", m, err)
		return nil, err
	}
	s.logger.Debugf("SecretText created: %v", m)
	return m, nil
}

func (s *server) addSecretFile(ctx context.Context, m *model.SecretFile, key, iv string) (*model.SecretFile, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if err := m.Encrypt(key, iv); err != nil {
		return nil, err
	}
	if err := s.store.SecretFile().Add(ctx, m); err != nil {
		s.logger.Errorf("Failed to add SecretFile %v: %v", m, err)
		return nil, err
	}
	s.logger.Debugf("SecretFile created: %v", m)
	return m, nil
}

func (s *server) updateLoginWithPassword(ctx context.Context, m *model.LoginWithPassword, key, iv string) (*model.LoginWithPassword, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if err := m.Encrypt(key, iv); err != nil {
		return nil, err
	}
	if err := s.store.LoginWithPassword().Update(ctx, m); err != nil {
		s.logger.Errorf("Failed to add LoginWithPassword %v: %v", m, err)
		return nil, err
	}
	s.logger.Debugf("LoginWithPassword updated: %v", m)
	return m, nil
}

func (s *server) updateCreditCard(ctx context.Context, m *model.CreditCard, key, iv string) (*model.CreditCard, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if err := m.Encrypt(key, iv); err != nil {
		return nil, err
	}
	if err := s.store.CreditCard().Update(ctx, m); err != nil {
		s.logger.Errorf("Failed to add CreditCard %v: %v", m, err)
		return nil, err
	}
	s.logger.Debugf("CreditCard updated: %v", m)
	return m, nil
}

func (s *server) updateSecretText(ctx context.Context, m *model.SecretText, key, iv string) (*model.SecretText, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if err := m.Encrypt(key, iv); err != nil {
		return nil, err
	}
	if err := s.store.SecretText().Update(ctx, m); err != nil {
		s.logger.Errorf("Failed to add SecretText %v: %v", m, err)
		return nil, err
	}
	s.logger.Debugf("SecretText updated: %v", m)
	return m, nil
}

func (s *server) updateSecretFile(ctx context.Context, m *model.SecretFile, key, iv string) (*model.SecretFile, error) {
	storageFile, err := s.store.SecretFile().GetByID(ctx, m.ID, m.UserID)
	if err != nil {
		s.logger.Errorf("Unable to find SecretFile in db %v: %v", m, err)
		return nil, err
	}
	m.Path = storageFile.Path
	if err := m.Validate(); err != nil {
		s.logger.Errorf("SecretFile validation failed %v: %v", m, err)
		return nil, err
	}
	if err := s.store.SecretFile().Update(ctx, m); err != nil {
		s.logger.Errorf("Failed to add SecretFile %v: %v", m, err)
		return nil, err
	}
	s.logger.Debugf("SecretFile updated: %v", m)
	return m, nil
}

func (s *server) deleteLoginWithPassword(ctx context.Context, id, userID int) error {
	if err := s.store.LoginWithPassword().Delete(ctx, id, userID); err != nil {
		return err
	}
	return nil
}

func (s *server) deleteCreditCard(ctx context.Context, id, userID int) error {
	if err := s.store.CreditCard().Delete(ctx, id, userID); err != nil {
		return err
	}
	return nil
}

func (s *server) deleteSecretText(ctx context.Context, id, userID int) error {
	if err := s.store.SecretText().Delete(ctx, id, userID); err != nil {
		return err
	}
	return nil
}

func (s *server) deleteSecretFile(ctx context.Context, id, userID int, key, iv string) error {
	m, err := s.store.SecretFile().GetByID(ctx, id, userID)
	if err != nil {
		return err
	}
	if err := m.Decrypt(key, iv); err != nil {
		return err
	}
	if err := s.store.SecretFile().Delete(ctx, id, userID); err != nil {
		return err
	}
	if err := os.Remove(m.Path); err != nil {
		return err
	}
	return nil
}

func (s *server) getLoginWithPassword(ctx context.Context, id, userID int, key, iv string) (*model.LoginWithPassword, error) {
	m, err := s.store.LoginWithPassword().GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if err := m.Decrypt(key, iv); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *server) getCreditCard(ctx context.Context, id, userID int, key, iv string) (*model.CreditCard, error) {
	m, err := s.store.CreditCard().GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if err := m.Decrypt(key, iv); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *server) getSecretText(ctx context.Context, id, userID int, key, iv string) (*model.SecretText, error) {
	m, err := s.store.SecretText().GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if err := m.Decrypt(key, iv); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *server) getSecretFile(ctx context.Context, id, userID int, key, iv string) (*model.SecretFile, error) {
	m, err := s.store.SecretFile().GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if err := m.Decrypt(key, iv); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *server) searchLoginWithPassword(ctx context.Context, name string, id int, key, iv string) ([]*model.LoginWithPassword, error) {
	m, err := s.store.LoginWithPassword().SearchByName(ctx, name, id)
	if err != nil {
		return nil, err
	}
	for _, i := range m {
		if err := i.Decrypt(key, iv); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (s *server) searchCreditCard(ctx context.Context, name string, id int, key, iv string) ([]*model.CreditCard, error) {
	m, err := s.store.CreditCard().SearchByName(ctx, name, id)
	if err != nil {
		return nil, err
	}
	for _, i := range m {
		if err := i.Decrypt(key, iv); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (s *server) searchSecretText(ctx context.Context, name string, id int, key, iv string) ([]*model.SecretText, error) {
	m, err := s.store.SecretText().SearchByName(ctx, name, id)
	if err != nil {
		return nil, err
	}
	for _, i := range m {
		if err := i.Decrypt(key, iv); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (s *server) searchSecretFile(ctx context.Context, name string, id int, key, iv string) ([]*model.SecretFile, error) {
	m, err := s.store.SecretFile().SearchByName(ctx, name, id)
	if err != nil {
		return nil, err
	}
	for _, i := range m {
		if err := i.Decrypt(key, iv); err != nil {
			return nil, err
		}
	}
	return m, nil
}
