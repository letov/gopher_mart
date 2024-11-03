package test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/db"
	"gopher_mart/internal/infrastructure/dto/request"
	"gopher_mart/internal/infrastructure/dto/response"
	"gopher_mart/internal/infrastructure/repo"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_SaveUser(t *testing.T) {
	type args struct {
		Login    string
		Password string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "login user test",
			args: args{
				Login:    "login",
				Password: "password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTest(t, func(cb *command.Bus, ur repo.User, er repo.Event, mux *chi.Mux, db *db.DB) {
				ctx := context.Background()
				_ = flushDB(ctx, db)

				data, _ := json.Marshal(request.SaveUser{
					Login:    tt.args.Login,
					Password: tt.args.Password,
					Name:     tt.args.Login,
				})
				req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(data))
				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)
				assert.Equal(t, rr.Code, http.StatusOK)

				req, _ = http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(data))
				mux.ServeHTTP(rr, req)
				assert.Equal(t, rr.Code, http.StatusBadRequest)
				assert.Equal(t, strings.TrimSpace(rr.Body.String()), command.ErrUserExist.Error())

				userExist, _ := er.HasEvent(ctx, tt.args.Login, domain.SaveUserAction, 0)
				assert.True(t, userExist)

				data, _ = json.Marshal(request.Login{
					Login:    "invalid",
					Password: "invalid",
				})
				req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(data))
				rr = httptest.NewRecorder()
				mux.ServeHTTP(rr, req)
				assert.Equal(t, rr.Code, http.StatusForbidden)
				assert.Equal(t, strings.TrimSpace(rr.Body.String()), command.ErrIncorrectLoginOrPassword.Error())

				time.Sleep(time.Millisecond * 200) // UserRepo update async by event

				data, _ = json.Marshal(request.Login{
					Login:    tt.args.Login,
					Password: tt.args.Password,
				})
				req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(data))
				rr = httptest.NewRecorder()
				mux.ServeHTTP(rr, req)
				assert.Equal(t, rr.Code, http.StatusOK)
				data, _ = io.ReadAll(rr.Body)
				_ = req.Body.Close()
				var jr response.Login
				_ = json.Unmarshal(data, &jr)
				assert.True(t, len(jr.Token) > 0)
			})
		})
	}
}
