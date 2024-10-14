package test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/dto/request"
	"gopher_mart/internal/infrastructure/dto/response"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_SaveUser(t *testing.T) {
	type args struct {
		reqDto request.SaveUser
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "save user test",
			args: args{
				reqDto: request.SaveUser{
					Login:    "login",
					Password: "password",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTest(t, func(cb *command.Bus, ur *UserRepo, er *EventRepo, mux *chi.Mux) {
				ctx := context.Background()
				data, _ := json.Marshal(tt.args.reqDto)

				req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(data))
				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				assert.True(t, er.HasEvent(ctx, tt.args.reqDto.Login, domain.SaveUserAction))
				time.Sleep(time.Millisecond * 200) // UserRepo update async by event
				assert.True(t, ur.HasUser(tt.args.reqDto.Login))

				req, _ = http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(data))
				rr = httptest.NewRecorder()
				mux.ServeHTTP(rr, req)
				assert.Equal(t, rr.Code, http.StatusBadRequest) // User already exist

				data, _ = json.Marshal(request.Login{
					Login:    tt.args.reqDto.Login,
					Password: "WRONG_PASS",
				})
				req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(data))
				rr = httptest.NewRecorder()
				mux.ServeHTTP(rr, req)
				assert.Equal(t, rr.Code, http.StatusForbidden)

				data, _ = json.Marshal(request.Login{
					Login:    tt.args.reqDto.Login,
					Password: tt.args.reqDto.Password,
				})
				req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(data))
				rr = httptest.NewRecorder()
				mux.ServeHTTP(rr, req)
				assert.Equal(t, rr.Code, http.StatusOK)
				var jr response.Login
				data, _ = io.ReadAll(rr.Body)
				_ = req.Body.Close()
				_ = json.Unmarshal(data, &jr)
				assert.True(t, len(jr.Token) > 0)
			})
		})
	}
}
