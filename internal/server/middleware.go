package server

import (
	"context"
	"fmt"
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
		id, ok := session.Values["authorization"]
		if !ok {
			fmt.Println("Unable to get session value")
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		u, err := s.store.User().FindByID(r.Context(), id.(int))
		if err != nil {
			fmt.Printf("Unable to find user: %v", err)
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

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
			"Request: %s %s %v %d %v",
			r.Method,
			r.RequestURI,
			time.Since(start),
			responseWriter.code,
			http.StatusText(responseWriter.code),
		)
	})
}
