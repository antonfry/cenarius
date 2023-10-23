package server

import (
	"cenarius/internal/model"
	"encoding/json"
	"net/http"
	"strconv"

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
	r.Get("/health", s.handleHealthCheck())

	r.Get("/loginwithpasswords", s.handleLoginWithPasswordSearch())
	r.Get("/loginwithpassword/{id}", s.handleLoginWithPasswordWithID())
	r.Get("/loginwithpassword/search/{name}", s.handleLoginWithPasswordSearch())
	r.Put("/loginwithpassword", s.handleLoginWithPasswordWithBody())
	r.Post("/loginwithpassword", s.handleLoginWithPasswordWithBody())
	r.Delete("/loginwithpassword/{id}", s.handleLoginWithPasswordWithID())

	r.Get("/creditcards", s.handleCreditCardSearch())
	r.Get("/creditcard/{id}", s.handleCreditCardWithID())
	r.Get("/creditcard/search/{name}", s.handleCreditCardSearch())
	r.Put("/creditcard", s.handleCreditCardWithBody())
	r.Post("/creditcard", s.handleCreditCardWithBody())
	r.Delete("/creditcard/{id}", s.handleCreditCardWithID())

	r.Get("/secrettexts", s.handleSecretTextSearch())
	r.Get("/secrettext/{id}", s.handleSecretTextWithID())
	r.Get("/secrettext/search/{name}", s.handleSecretTextSearch())
	r.Put("/secrettext", s.handleSecretTextWithBody())
	r.Post("/secrettext", s.handleSecretTextWithBody())
	r.Delete("/secrettext/{id}", s.handleSecretTextWithID())

	r.Get("/secretfiles", s.handleSecretFileSearch())
	r.Get("/secretfile/{id}", s.handleSecretFileWithID())
	r.Get("/secretfile/search/{name}", s.handleSecretFileSearch())
	r.Put("/secretfile", s.handleSecretFileWithBody())
	r.Post("/secretfile", s.handleSecretFileWithBody())
	r.Delete("/secretfile/{id}", s.handleSecretFileWithID())

	return r
}

func (s *server) saveSession(w http.ResponseWriter, r *http.Request, u *model.User) error {
	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		s.logger.Errorf("unable to get session %v", sessionName)
		return err
	}
	s.logger.Debugf("Saving session: %v", u.ID)
	session.Values["authorization"] = u.ID
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
			return
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
		s.logger.Info("handleUserLogin is working")
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

func (s *server) handleLoginWithPasswordWithBody() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.LoginWithPassword{}
		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			s.logger.Errorf("Unable to parse body in handleAddLoginWithPassword: %v", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		switch r.Method {
		case "PUT":
			if _, err := s.addLoginWithPassword(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		case "POST":
			if _, err := s.updateLoginWithPassword(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleLoginWithPasswordWithID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.LoginWithPassword{}
		var err error
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		m.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		switch r.Method {
		case "GET":
			if _, err := s.getLoginWithPassword(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		case "DELETE":
			if err := s.deleteLoginWithPassword(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleLoginWithPasswordSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKeyUser)
		userId := user.(*model.User).ID
		name := chi.URLParam(r, "name")
		if _, err := s.searchLoginWithPassword(r.Context(), name, userId); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}
}

func (s *server) handleCreditCardWithBody() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.CreditCard{}
		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			s.logger.Errorf("Unable to parse body in handleAddCreditCard: %v", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		switch r.Method {
		case "PUT":
			if _, err := s.addCreditCard(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		case "POST":
			if _, err := s.updateCreditCard(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleCreditCardWithID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.CreditCard{}
		var err error
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		m.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		switch r.Method {
		case "GET":
			if _, err := s.getCreditCard(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		case "DELETE":
			if err := s.deleteCreditCard(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleCreditCardSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKeyUser)
		userId := user.(*model.User).ID
		name := chi.URLParam(r, "name")
		if _, err := s.searchCreditCard(r.Context(), name, userId); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}
}

func (s *server) handleSecretTextWithBody() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.SecretText{}
		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			s.logger.Errorf("Unable to parse body in handleAddSecretText: %v", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		switch r.Method {
		case "PUT":
			if _, err := s.addSecretText(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		case "POST":
			if _, err := s.updateSecretText(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleSecretTextWithID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.SecretText{}
		var err error
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		m.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		switch r.Method {
		case "GET":
			if _, err := s.getSecretText(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		case "DELETE":
			if err := s.deleteSecretText(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleSecretTextSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKeyUser)
		userId := user.(*model.User).ID
		name := chi.URLParam(r, "name")
		if _, err := s.searchSecretText(r.Context(), name, userId); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}
}

func (s *server) handleSecretFileWithBody() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.SecretFile{}
		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			s.logger.Errorf("Unable to parse body in handleAddSecretFile: %v", err)
			s.error(w, r, http.StatusBadRequest, err)
		}
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		switch r.Method {
		case "PUT":
			if _, err := s.addSecretFile(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		case "POST":
			if _, err := s.updateSecretFile(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleSecretFileWithID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.SecretFile{}
		var err error
		user := r.Context().Value(ctxKeyUser)
		m.UserId = user.(*model.User).ID
		m.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		switch r.Method {
		case "GET":
			if m, err = s.getSecretFile(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
			uploadWS(w, r, m)
			return
		case "DELETE":
			if err := s.deleteSecretFile(r.Context(), m); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
			s.respond(w, r, http.StatusOK, m)
		}
	}
}

func (s *server) handleSecretFileSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKeyUser)
		userId := user.(*model.User).ID
		name := chi.URLParam(r, "name")
		if _, err := s.searchSecretFile(r.Context(), name, userId); err != nil {
			s.logger.Errorf("func handleSecretFileSearch: %s", err.Error())
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}
}
