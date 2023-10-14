package server

import (
	"cenarius/internal/model"
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
	s.router.Post("/api/v1/user/login", s.handleUserLogin())
	s.router.Post("/api/v1/user/register", s.handleUserRegister())
	s.router.Get("/ping", s.handleHealthCheck())

	s.router.Mount("/api/v1/private", s.privateRouter())
	s.HTTPServer.Handler = s.router
}

func (s *server) privateRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(s.authenticateUser)
	r.Get("/p", s.handleHealthCheck())

	r.Get("/loginwithpassword/{id}", s.handleHealthCheck())
	r.Get("/loginwithpassword/search/{name}", s.handleHealthCheck())
	r.Put("/loginwithpassword", s.handleLoginWithPassword())
	r.Post("/loginwithpassword", s.handleLoginWithPassword())
	r.Delete("/loginwithpassword", s.handleLoginWithPassword())

	r.Get("/creditcard/{id}", s.handleHealthCheck())
	r.Put("/creditcard", s.handleAddCreditCard())
	r.Post("/creditcard", s.handleHealthCheck())
	r.Delete("/creditcard", s.handleDeleteCreditCard())

	r.Get("/secrettext/{id}", s.handleHealthCheck())
	r.Put("/secrettext", s.handleAddSecretText())
	r.Post("/secrettext", s.handleHealthCheck())
	r.Delete("/secrettext", s.handleDeleteSecretText())

	r.Get("/secretbinary/{id}", s.handleHealthCheck())
	r.Put("/secretbinary", s.handleAddSecretBinary())
	r.Post("/secretbinary", s.handleHealthCheck())
	r.Delete("/secretbinary", s.handleDeleteSecretBinary())
	return r
}

func (s *server) saveSession(w http.ResponseWriter, r *http.Request, u *model.User) error {
	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		s.logger.Errorf("unable to get session %v", sessionName)
		return err
	}
	s.logger.Infof("Saving session: %v", u.ID)
	session.Values["authorization"] = u.ID
	s.logger.Infof("Saving session: %v", session.Values["authorization"])
	if err := s.sessionStore.Save(r, w, session); err != nil {
		s.logger.Errorf("unable to save session for user %v", u)
		return err
	}
	return nil
}

func (s *server) handleUserRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &model.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			s.logger.Errorf("Unable to parse body: %v", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		u, code, err := s.userRegister(r.Context(), u)
		if err != nil {
			s.error(w, r, code, err)
			return
		}
		err = s.saveSession(w, r, u)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, u)
	}
}

func (s *server) handleUserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &model.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			s.logger.Errorf("Unable to parse body: %v", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		u, code, err := s.userLogin(r.Context(), u)
		if err != nil {
			s.error(w, r, code, err)
			return
		}
		err = s.saveSession(w, r, u)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, u)
	}
}

func (s *server) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusNoContent, nil)
	}
}

func (s *server) handleLoginWithPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.LoginWithPassword{}
		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			s.logger.Errorf("Unable to parse body in handleAddLoginWithPassword: %v", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		switch r.Method {
		case "GET":
			if _, err := s.addLoginWithPassword(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		case "PUT":
			if _, err := s.addLoginWithPassword(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		case "POST":
			if _, err := s.updateLoginWithPassword(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		case "DELETE":
			if err := s.deleteLoginWithPassword(r.Context(), m.ID); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		}

		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleAddCreditCard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.CreditCard{}
		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			s.logger.Errorf("Unable to parse body in handleAddCreditCard: %v", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		if _, err := s.addCreditCard(r.Context(), m); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusNoContent, nil)
	}
}

func (s *server) handleAddSecretText() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.SecretText{}
		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			s.logger.Errorf("Unable to parse body in handleAddSecretText: %v", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		if _, err := s.addSecretText(r.Context(), m); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusNoContent, nil)
	}
}

func (s *server) handleAddSecretBinary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.SecretBinary{}
		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			s.logger.Errorf("Unable to parse body in handleAddSecretBinary: %v", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		if _, err := s.addSecretBinary(r.Context(), m); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusNoContent, nil)
	}
}

func (s *server) handleDeleteLoginWithPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.LoginWithPassword{}
		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			s.logger.Errorf("Unable to parse body in handleDeleteLoginWithPassword: %v", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		if err := s.deleteLoginWithPassword(r.Context(), m.ID); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusNoContent, m)
	}
}

func (s *server) handleDeleteCreditCard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.LoginWithPassword{}
		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			s.logger.Errorf("Unable to parse body in handleDeleteLoginWithPassword: %v", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		if err := s.deleteCreditCard(r.Context(), m.ID); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusNoContent, m)
	}
}

func (s *server) handleDeleteSecretText() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.SecretText{}
		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			s.logger.Errorf("Unable to parse body in handleDeleteLoginWithPassword: %v", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		if err := s.deleteSecretText(r.Context(), m.ID); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusNoContent, m)
	}
}

func (s *server) handleDeleteSecretBinary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.SecretBinary{}
		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			s.logger.Errorf("Unable to parse body in handleDeleteLoginWithPassword: %v", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		if err := s.deleteSecretBinary(r.Context(), m.ID); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusNoContent, m)
	}
}
