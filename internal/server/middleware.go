package server

import (
	"cenarius/internal/store"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (s *server) setContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.logger.Errorf("authenticateUser: unable to get session %v", sessionName)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.logger.Infof("server.authenticateUser session.Values: %v", session.Values)
		c, ok := session.Values[sessionName]
		if !ok {
			s.logger.Errorf("Unable to get session value")
			s.error(w, r, http.StatusUnauthorized, store.ErrNotAuthenticated)
			return
		}

		sessionData, ok := c.(cenariusSession)
		if !ok {
			s.logger.Errorf("Unable to determine session value")
			s.error(w, r, http.StatusUnauthorized, store.ErrNotAuthenticated)
			return
		}
		s.logger.Infof("Got cookie: %v", c)
		idInt, ok := c
		if !ok {
			s.error(w, r, http.StatusBadRequest, store.ErrNotAuthenticated)
			return
		}
		u, err := s.store.User().FindByID(r.Context(), idInt)
		if err != nil {
			s.logger.Errorf("Unable to find user from session in store : %v", err)
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}
		s.logger.Info("server.authenticateUser ok")
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger := s.logger.WithFields(logrus.Fields{
			"ip": r.RemoteAddr,
			"id": r.Context().Value(ctxKeyRequestID),
		})
		responseWriter := &responseWriter{w, 0}
		next.ServeHTTP(responseWriter, r)
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
