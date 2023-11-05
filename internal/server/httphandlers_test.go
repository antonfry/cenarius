package server

import (
	"bytes"
	"cenarius/internal/model"
	"cenarius/internal/store/sqlstore"
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("CENARIUS_DATABASEDSN")
	if databaseURL == "" {
		databaseURL = "host=localhost dbname=cenarius_test sslmode=disable"
	}
	os.Exit(m.Run())
}

func Test_server_handleUserRegister(t *testing.T) {
	_, teardown := sqlstore.TestStore(t, databaseURL)
	defer teardown("users")
	tests := []struct {
		name string
		m    *model.User
		want int
	}{
		{
			name: "Valid",
			m:    &model.User{Login: "Valid", Password: "testpassword"},
			want: 200,
		},
		{
			name: "AlreadyExist",
			m:    &model.User{Login: "Valid", Password: "testpassword"},
			want: 409,
		},
		{
			name: "InValid",
			m:    &model.User{Login: "", Password: "testpassword"},
			want: 400,
		},
	}
	conf := NewConfig()
	conf.DatabaseDsn = databaseURL
	s := NewServer(conf)
	handler := http.HandlerFunc(s.handleUserRegister())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.m)
			var buf bytes.Buffer
			_, _ = io.WriteString(&buf, string(jsonData))
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/api/v1/user/register", &buf)
			if err != nil {
				t.Errorf("http.NewRequest error = %v", err)
			}
			defer req.Body.Close()
			handler.ServeHTTP(rec, req)
			assert.Equal(t, tt.want, rec.Result().StatusCode)
		})
	}
}

func Test_server_handleHealthCheck(t *testing.T) {
	conf := NewConfig()
	s := NewServer(conf)
	handler := http.HandlerFunc(s.handleHealthCheck())
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Result().StatusCode)
}

func Test_server_handleLoginWithPasswordWithBody(t *testing.T) {
	conf := NewConfig()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	u := &model.User{Login: "Valid", EncryptedPassword: "testpasswordtestpasswordtestpass", ID: r.Intn(1000-10) + 1}
	s := NewServer(conf)
	handler := http.HandlerFunc(s.handleLoginWithPasswordWithBody())
	tests := []struct {
		name   string
		method string
		uri    string
		m      *model.LoginWithPassword
		ctx    context.Context
		want   int
	}{
		{
			name:   "ValidPost",
			method: http.MethodPost,
			m:      &model.LoginWithPassword{Login: "TestValidLoginPost", Password: "TestValidPasswordPost"},
			ctx:    context.WithValue(context.Background(), ctxKeyUser, u),
			want:   http.StatusOK,
		},
		{
			name:   "ValidPut",
			method: http.MethodPut,
			m:      &model.LoginWithPassword{Login: "TestValidLoginPut", Password: "TestValidPasswordPut"},
			ctx:    context.WithValue(context.Background(), ctxKeyUser, u),
			want:   http.StatusOK,
		},
		{
			name:   "InValidPost",
			method: http.MethodPost,
			m:      &model.LoginWithPassword{Login: "", Password: ""},
			ctx:    context.WithValue(context.Background(), ctxKeyUser, u),
			want:   http.StatusInternalServerError,
		},
		{
			name:   "InValidPut",
			method: http.MethodPut,
			m:      &model.LoginWithPassword{Login: "", Password: ""},
			ctx:    context.WithValue(context.Background(), ctxKeyUser, u),
			want:   http.StatusInternalServerError,
		},
		{
			name:   "InValidPostContext",
			method: http.MethodPost,
			m:      &model.LoginWithPassword{Login: "InValidPostContext", Password: "InValidPostContextpass"},
			ctx:    context.Background(),
			want:   http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			jsonData, _ := json.Marshal(tt.m)
			var buf bytes.Buffer
			_, _ = io.WriteString(&buf, string(jsonData))
			req, err := http.NewRequest(tt.method, "/", &buf)
			if err != nil {
				t.Error(err)
			}
			defer req.Body.Close()
			handler.ServeHTTP(rec, req.WithContext(tt.ctx))
			assert.Equal(t, tt.want, rec.Result().StatusCode)
		})
	}
}
