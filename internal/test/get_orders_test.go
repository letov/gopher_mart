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
	"testing"
	"time"
)

func Test_GetOrder(t *testing.T) {
	type args struct {
		Login    string
		Password string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "get orders test",
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

				createOrder(response.AccrualOrder{
					OrderID: "17893729974",
					Status:  response.RegisteredStatus,
					Accrual: 312,
				}, token, mux, h)
				createOrder(response.AccrualOrder{
					OrderID: "123456789031",
					Status:  response.InvalidStatus,
					Accrual: 412,
				}, token, mux, h)

				req, _ := http.NewRequest("GET", "/api/user/orders", nil)
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				data := rr.Body.String()
				var dto []response.Order
				_ = json.Unmarshal([]byte(data), &dto)

				assert.Equal(t, len(dto), 2)
				assert.Equal(t, dto[0].OrderID, "123456789031")
				assert.Equal(t, dto[1].OrderID, "17893729974")
				assert.Equal(t, dto[1].Status, domain.NewStatus)
			})
		})
	}
}

func createOrder(o response.AccrualOrder, token string, mux *chi.Mux, h *HttpClient) {
	h.SetResponse(o)

	data, _ := json.Marshal(request.SaveOrder{
		OrderID: o.OrderID,
	})
	req, _ := http.NewRequest("POST", "/api/user/orders", bytes.NewBuffer(data))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	time.Sleep(time.Millisecond * 200) // ждем асинхронной обработки событий
}
