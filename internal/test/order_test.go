package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/db"
	"gopher_mart/internal/infrastructure/dto/request"
	"gopher_mart/internal/infrastructure/dto/response"
	"gopher_mart/internal/infrastructure/queue"
	"gopher_mart/internal/infrastructure/repo"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_RequestAccrual(t *testing.T) {
	type args struct {
		Login    string
		Password string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "request order accrual user test",
			args: args{
				Login:    "login",
				Password: "password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTest(t, func(
				cb *command.Bus,
				ur repo.User,
				er repo.Event,
				or repo.Order,
				mux *chi.Mux,
				db *db.DB,
				q *queue.Rabbit,
				h *HttpClient,
			) {
				ctx := context.Background()
				_ = flushDB(ctx, db)

				_, _ = cb.Execute(command.SaveUser{
					Ctx: ctx,
					Request: request.SaveUser{
						Login:    tt.args.Login,
						Password: tt.args.Password,
						Name:     tt.args.Login,
					},
				})

				rl, _ := cb.Execute(command.Login{
					Ctx: ctx,
					Request: request.Login{
						Login:    tt.args.Login,
						Password: tt.args.Password,
					},
				})

				token := rl.(response.Login).Token

				invalidOrderId := "123"
				data, _ := json.Marshal(request.RequestAccrual{
					OrderID: invalidOrderId,
				})
				req, _ := http.NewRequest("POST", "/api/user/orders", bytes.NewBuffer(data))
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)
				assert.Equal(t, rr.Code, http.StatusBadRequest)
				assert.Equal(t, strings.TrimSpace(rr.Body.String()), command.ErrInvalidOrderId.Error())

				validOrderId := "17893729974"
				accrual := int64(1000)

				h.SetResponse(response.OrderAccrual{
					OrderID: validOrderId,
					Status:  response.RegisteredStatus,
					Accrual: accrual,
				})

				data, _ = json.Marshal(request.RequestAccrual{
					OrderID: validOrderId,
				})
				req, _ = http.NewRequest("POST", "/api/user/orders", bytes.NewBuffer(data))
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
				rr = httptest.NewRecorder()
				mux.ServeHTTP(rr, req)
				assert.Equal(t, rr.Code, http.StatusOK)

				time.Sleep(time.Millisecond * 200) // waiting queue handler processed

				o, _ := or.Get(ctx, validOrderId)
				assert.Equal(t, o.Status, domain.NewStatus)
				assert.Equal(t, o.Accrual, accrual)
			})
		})
	}
}
