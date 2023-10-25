package server

import (
	"cenarius/internal/model"
	"cenarius/internal/store"
	"cenarius/internal/store/sqlstore"
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type ctxKey int8

const (
	AuthHeader        = "X-Cenarius-Token"
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

// server server main struct
type server struct {
	config        *Config
	logger        *logrus.Logger
	HTTPServer    *http.Server
	router        *chi.Mux
	GRPCServer    *grpc.Server
	store         store.Store
	allowedSubnet *net.IPNet
}

// NewServer returns new server object
func NewServer(config *Config) *server {
	s := &server{
		config:     config,
		logger:     logrus.New(),
		HTTPServer: &http.Server{Addr: config.Bind},
	}
	s.configureLogger()
	s.configureStore()
	s.configureTrustedSubnets()

	return s
}

// StartHTTPServer starts GRPC Server
func (s *server) StartGRPCServer() {
	s.logger.Infof("Starting Grpc server with config: %v\n", s.config)
	listen, err := net.Listen("tcp", s.config.Bind)
	if err != nil {
		s.logger.Fatal(err)
	}
	s.GRPCServer = grpc.NewServer(grpc.UnaryInterceptor(s.unaryInterceptor))
	if err := s.GRPCServer.Serve(listen); err != nil {
		s.logger.Fatal(err)
	}
	s.logger.Infof("GRPC server stopped with config: %v\n", s.config)
}

// StopGRPCServer stops grpc server
func (s *server) StopGRPCServer() {
	s.GRPCServer.GracefulStop()
}

// StartHTTPServer starts GRPC Server
func (s *server) StartHTTPServer() {
	s.logger.Infof("Starting HTTP server with config: %v\n", s.config)
	s.router = chi.NewRouter()
	s.configureRouter()
	if err := s.HTTPServer.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Fatal(err)
	}
	s.logger.Infof("HTTP server stopped with config: %v\n", s.config)
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
	switch {
	case s.config.Mode == "GRPC":
		s.StartGRPCServer()
	case s.config.Mode == "HTTP":
		s.StartHTTPServer()
	default:
		s.logger.Fatalf("Unknow node %s", s.config.Mode)
		return nil
	}
	return nil
}

// Shutdown shutdowns the server
func (s *server) Shutdown() {
	s.logger.Info("Shuting down...")
	switch {
	case s.config.Mode == "GRPC":
		s.StopGRPCServer()
	case s.config.Mode == "HTTP":
		s.StopHTTPServer()
	default:
		s.logger.Fatalf("Unknow node %s", s.config.Mode)
	}
	s.store.Close()
	s.logger.Info("Done ShutDown")
}

// configureLogger configures logger
func (s *server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

// configureStore configures store
func (s *server) configureStore() {
	conn, err := sqlstore.NewPGConn(s.config.DatabaseDsn)
	if err != nil {
		s.logger.Errorf("Unable to connect to the database with: %v", s.config.DatabaseDsn)
		s.logger.Fatal(err)
	}
	s.store = sqlstore.NewStore(conn)
}

// configureTrustedSubnets configures trusted subnets
func (s *server) configureTrustedSubnets() {
	if s.config.TrustedSubnet != "" {
		_, subnet, err := net.ParseCIDR(s.config.TrustedSubnet)
		if err != nil {
			log.Fatalf("Can't parse subnet from %s %v", s.config.TrustedSubnet, err)
		}
		s.allowedSubnet = subnet
	}
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
	s.logger.Infof("User created: %v", u)
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
	u.Sanitaze()
	return u, nil
}

func (s *server) addLoginWithPassword(ctx context.Context, m *model.LoginWithPassword) (*model.LoginWithPassword, error) {
	if err := s.store.LoginWithPassword().Add(ctx, m); err != nil {
		s.logger.Errorf("Failed to add LoginWithPassword %v: %v", m, err)
		return nil, err
	}
	s.logger.Infof("LoginWithPassword created: %v", m)
	return m, nil
}

func (s *server) addCreditCard(ctx context.Context, m *model.CreditCard) (*model.CreditCard, error) {
	if err := s.store.CreditCard().Add(ctx, m); err != nil {
		s.logger.Errorf("Failed to add CreditCard %v: %v", m, err)
		return nil, err
	}
	s.logger.Infof("CreditCard created: %v", m)
	return m, nil
}

func (s *server) addSecretText(ctx context.Context, m *model.SecretText) (*model.SecretText, error) {
	if err := s.store.SecretText().Add(ctx, m); err != nil {
		s.logger.Errorf("Failed to add SecretText %v: %v", m, err)
		return nil, err
	}
	s.logger.Infof("SecretText created: %v", m)
	return m, nil
}

func (s *server) addSecretFile(ctx context.Context, m *model.SecretFile) (*model.SecretFile, error) {

	if err := s.store.SecretFile().Add(ctx, m); err != nil {
		s.logger.Errorf("Failed to add SecretFile %v: %v", m, err)
		return nil, err
	}
	s.logger.Infof("SecretFile created: %v", m)
	return m, nil
}

func (s *server) updateLoginWithPassword(ctx context.Context, m *model.LoginWithPassword) (*model.LoginWithPassword, error) {
	if err := s.store.LoginWithPassword().Update(ctx, m); err != nil {
		s.logger.Errorf("Failed to add LoginWithPassword %v: %v", m, err)
		return nil, err
	}
	s.logger.Infof("LoginWithPassword created: %v", m)
	return m, nil
}

func (s *server) updateCreditCard(ctx context.Context, m *model.CreditCard) (*model.CreditCard, error) {
	if err := s.store.CreditCard().Update(ctx, m); err != nil {
		s.logger.Errorf("Failed to add CreditCard %v: %v", m, err)
		return nil, err
	}
	s.logger.Infof("CreditCard created: %v", m)
	return m, nil
}

func (s *server) updateSecretText(ctx context.Context, m *model.SecretText) (*model.SecretText, error) {
	if err := s.store.SecretText().Update(ctx, m); err != nil {
		s.logger.Errorf("Failed to add SecretText %v: %v", m, err)
		return nil, err
	}
	s.logger.Infof("SecretText created: %v", m)
	return m, nil
}

func (s *server) updateSecretFile(ctx context.Context, m *model.SecretFile) (*model.SecretFile, error) {
	if err := s.store.SecretFile().Update(ctx, m); err != nil {
		s.logger.Errorf("Failed to add SecretFile %v: %v", m, err)
		return nil, err
	}
	s.logger.Infof("SecretFile created: %v", m)
	return m, nil
}

func (s *server) deleteLoginWithPassword(ctx context.Context, m *model.LoginWithPassword) error {
	if err := s.store.LoginWithPassword().Delete(ctx, m); err != nil {
		return err
	}
	return nil
}

func (s *server) deleteCreditCard(ctx context.Context, m *model.CreditCard) error {
	if err := s.store.CreditCard().Delete(ctx, m); err != nil {
		return err
	}
	return nil
}

func (s *server) deleteSecretText(ctx context.Context, m *model.SecretText) error {
	if err := s.store.SecretText().Delete(ctx, m); err != nil {
		return err
	}
	return nil
}

func (s *server) deleteSecretFile(ctx context.Context, m *model.SecretFile) error {
	if err := s.store.SecretFile().Delete(ctx, m); err != nil {
		return err
	}
	return nil
}

func (s *server) getLoginWithPassword(ctx context.Context, m *model.LoginWithPassword) (*model.LoginWithPassword, error) {
	m, err := s.store.LoginWithPassword().GetByID(ctx, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *server) getCreditCard(ctx context.Context, m *model.CreditCard) (*model.CreditCard, error) {
	m, err := s.store.CreditCard().GetByID(ctx, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *server) getSecretText(ctx context.Context, m *model.SecretText) (*model.SecretText, error) {
	m, err := s.store.SecretText().GetByID(ctx, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *server) getSecretFile(ctx context.Context, m *model.SecretFile) (*model.SecretFile, error) {
	m, err := s.store.SecretFile().GetByID(ctx, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *server) searchLoginWithPassword(ctx context.Context, name string, id int) ([]*model.LoginWithPassword, error) {
	m, err := s.store.LoginWithPassword().SearchByName(ctx, name, id)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *server) searchCreditCard(ctx context.Context, name string, id int) ([]*model.CreditCard, error) {
	m, err := s.store.CreditCard().SearchByName(ctx, name, id)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *server) searchSecretText(ctx context.Context, name string, id int) ([]*model.SecretText, error) {
	m, err := s.store.SecretText().SearchByName(ctx, name, id)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *server) searchSecretFile(ctx context.Context, name string, id int) ([]*model.SecretFile, error) {
	m, err := s.store.SecretFile().SearchByName(ctx, name, id)
	if err != nil {
		return nil, err
	}
	return m, nil
}
