package server

import (
	"cenarius/internal/model"
	"cenarius/internal/store"
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (s *server) setContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debug("server.setContentType is working")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debug("server.authenticateUser is working")
		h := r.Header.Get(AuthHeader)
		if h == "" {
			s.logger.Error("Unable to get auth header")
			s.error(w, r, http.StatusUnauthorized, store.ErrNotAuthenticated)
			return
		}
		data, err := base64.StdEncoding.DecodeString(h)
		if err != nil {
			s.logger.Errorf("Unable to decode header %s: %s", h, err.Error())
			s.error(w, r, http.StatusUnauthorized, store.ErrNotAuthenticated)
			return
		}
		loginAndPassword := strings.Fields(string(data))
		if len(loginAndPassword) != 2 {
			s.logger.Errorf("Bad header: %s", h)
			s.error(w, r, http.StatusUnauthorized, store.ErrNotAuthenticated)
			return
		}
		u := &model.User{Login: loginAndPassword[0], Password: loginAndPassword[1]}
		u, err = s.userLogin(r.Context(), u)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, store.ErrNotAuthenticated)
			return
		}
		s.logger.Debugf("server.authenticateUser ok: %s", u.Login)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debug("server.setRequestID is working")
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debug("server.logRequest is working")
		start := time.Now()
		logger := s.logger.WithFields(logrus.Fields{
			"ip": r.RemoteAddr,
			"id": r.Context().Value(ctxKeyRequestID),
		})
		responseWriter := &responseWriter{w, 0}
		next.ServeHTTP(responseWriter, r)
		s.logger.Debug("server.logRequest is logging")
		logger.Infof(
			"Request: %s %s %v %d %v %v",
			r.Method,
			r.RequestURI,
			time.Since(start),
			responseWriter.code,
			http.StatusText(responseWriter.code),
			r.UserAgent(),
		)
	})
}
