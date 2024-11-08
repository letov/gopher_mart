package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/application/service"
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

func Test_SaveOrder(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "order test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTest(t, func(
				cb *command.Bus,
				ur repo.User,
				er repo.Event,
				or repo.Order,
				opr repo.Operation,
				mux *chi.Mux,
				db *db.DB,
				q *queue.Rabbit,
				h *HttpClient,
				as *service.Accrual,
			) {
				ctx := context.Background()
				_ = flushDB(ctx, db)

				login := "login"
				password := "password"

				token := getToken(ctx, cb, request.Login{
					Login:    login,
					Password: password,
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

				h.SetResponse(response.AccrualOrder{
					OrderID: validOrderId,
					Status:  response.ProcessedStatus,
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
				assert.Equal(t, http.StatusOK, rr.Code) // номер заказа уже был загружен этим пользователем;

				newToken := getToken(ctx, cb, request.Login{
					Login:    "NewLogin",
					Password: "NewPassword",
				})
				data, _ = json.Marshal(request.SaveOrder{
					OrderID: validOrderId,
				})
				req, _ = http.NewRequest("POST", "/api/user/orders", bytes.NewBuffer(data))
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", newToken))
				rr = httptest.NewRecorder()
				mux.ServeHTTP(rr, req)
				assert.Equal(t, http.StatusConflict, rr.Code) // номер заказа уже был загружен другим пользователем;

				o, _ := or.Get(ctx, validOrderId)
				assert.Equal(t, domain.ProcessedStatus, o.Status)
				assert.Equal(t, accrual, o.Accrual)

				req, _ = http.NewRequest("GET", "/api/user/orders", nil)
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
				rr = httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				body := rr.Body.String()
				var dto []response.Order
				_ = json.Unmarshal([]byte(body), &dto)

				assert.Equal(t, validOrderId, dto[0].OrderID)
				assert.Equal(t, o.Accrual, dto[0].Accrual)

				op, _ := opr.GetByUserId(ctx, o.UserID)
				assert.Equal(t, domain.AddedStatus, op[0].Status)
				assert.Equal(t, accrual, op[0].Sum)
			})
		})
	}
}
