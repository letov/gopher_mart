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
			name: "load order test",
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

				token := getToken(ctx, cb, request.Login{
					Login:    tt.args.Login,
					Password: tt.args.Password,
				})

				invalidOrderId := "123"
				data, _ := json.Marshal(request.SaveOrder{
					OrderID: invalidOrderId,
				})
				req, _ := http.NewRequest("POST", "/api/user/orders", bytes.NewBuffer(data))
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)
				assert.Equal(t, rr.Code, http.StatusUnprocessableEntity) // неверный формат номера заказа;
				assert.Equal(t, strings.TrimSpace(rr.Body.String()), command.ErrInvalidOrderId.Error())

				validOrderId := "17893729974"
				accrual := int64(1000)

				h.SetResponse(response.Order{
					OrderID: validOrderId,
					Status:  response.RegisteredStatus,
					Accrual: accrual,
				})

				data, _ = json.Marshal(request.SaveOrder{
					OrderID: validOrderId,
				})
				req, _ = http.NewRequest("POST", "/api/user/orders", bytes.NewBuffer(data))
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
				rr = httptest.NewRecorder()
				mux.ServeHTTP(rr, req)
				assert.Equal(t, rr.Code, http.StatusAccepted) // новый номер заказа принят в обработку;

				data, _ = json.Marshal(request.SaveOrder{
					OrderID: validOrderId,
				})
				req, _ = http.NewRequest("POST", "/api/user/orders", bytes.NewBuffer(data))
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
				rr = httptest.NewRecorder()
				mux.ServeHTTP(rr, req)
				assert.Equal(t, rr.Code, http.StatusOK) // номер заказа уже был загружен этим пользователем;

				newToken := getToken(ctx, cb, request.Login{
					Login:    "NewLogin",
					Password: tt.args.Password,
				})
				data, _ = json.Marshal(request.SaveOrder{
					OrderID: validOrderId,
				})
				req, _ = http.NewRequest("POST", "/api/user/orders", bytes.NewBuffer(data))
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", newToken))
				rr = httptest.NewRecorder()
				mux.ServeHTTP(rr, req)
				assert.Equal(t, rr.Code, http.StatusConflict) // номер заказа уже был загружен другим пользователем;
				return

				o, _ := or.Get(ctx, validOrderId)
				assert.Equal(t, o.Status, domain.NewStatus)
				assert.Equal(t, o.Accrual, accrual)
			})
		})
	}
}
