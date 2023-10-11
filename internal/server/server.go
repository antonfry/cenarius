package server

import (
	"cenarius/internal/store"
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type ctxKey int8

const (
	sessionName        = "cenarius"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

type server struct {
	config        *Config
	logger        *logrus.Logger
	HTTPServer    *http.Server
	router        *chi.Mux
	sessionStore  sessions.Store
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
	s.GRPCServer = grpc.NewServer(grpc.UnaryInterceptor(s.unaryInterceptor))
	if config.TrustedSubnet != "" {
		_, subnet, err := net.ParseCIDR(config.TrustedSubnet)
		if err != nil {
			log.Fatalf("Can't parse subnet from %s %v", config.TrustedSubnet, err)
		}
		s.allowedSubnet = subnet
	}
	return s
}

// StartHTTPServer starts GRPC Server
func (s *server) StartGRPCServer() {
	s.logger.Infof("Starting Grpc server with config: %v\n", s.config)
	listen, err := net.Listen("tcp", s.config.Bind)
	if err != nil {
		s.logger.Fatal(err)
	}
	if err := s.GRPCServer.Serve(listen); err != nil {
		s.logger.Fatal(err)
	}
	s.logger.Infof("GRPC server stopped with config: %v\n", s.config)
}

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
	s.configureLogger()
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
