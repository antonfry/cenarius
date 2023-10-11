package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(gzipHandle)
	s.router.Use(s.setContentType)
	// s.router.Post("/api/user/login", s.handleUserLogin())
	// s.router.Post("/api/user/register", s.handleUserRegister())
	// s.router.Get("/ping", s.handleHealthCheck())

	s.router.Mount("/api", s.privateRouter())
	s.HTTPServer.Handler = s.router
}

func (s *server) privateRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(s.authenticateUser)
	// r.Get("/user/orders", s.handleUserOrdersGet())
	// r.Get("/user/balance", s.handleUserBalance())
	// r.Get("/user/withdrawals", s.handleUserWithdrawals())
	// r.Post("/user/orders", s.handleUserOrdersPost())
	// r.Post("/user/balance/withdraw", s.handleAddWithdrawal())
	return r
}
