package server

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		if r.Method == http.MethodGet || r.Method == http.MethodDelete {
			next.ServeHTTP(w, r)
			return
		}
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			r.Body, err = gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()
		}
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gzWriter, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			if _, err := io.WriteString(w, err.Error()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		defer gzWriter.Close()
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", r.Header.Get("Accept"))
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gzWriter}, r)
	})
}
