package test

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/infrastructure/dto"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_SaveUser(t *testing.T) {
	type args struct {
		reqDto dto.SaveUser
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "save user test",
			args: args{
				reqDto: dto.SaveUser{
					Login:    "login",
					Password: "password",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intTest(t, func(cb *command.Bus, ur *UserRepo, er *EventRepo, mux *chi.Mux) {
				data, _ := json.Marshal(tt.args.reqDto)

				req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(data))
				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				assert.True(t, er.HasRootId(tt.args.reqDto.Login))
				time.Sleep(time.Millisecond * 200) // UserRepo update async by event
				assert.True(t, ur.HasLogin(tt.args.reqDto.Login))

				req, _ = http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(data))
				rr = httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				assert.Equal(t, rr.Code, http.StatusBadRequest) // User already exist
			})
		})
	}
}
