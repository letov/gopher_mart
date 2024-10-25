package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/infrastructure/dto/response"
	"io"
	"net/http"
)

var (
	ErrNoOKResponse = errors.New("no ok response")
)

type HttpClient struct {
	config *config.Config
}

func NewHttpClient(
	config *config.Config,
) *HttpClient {
	return &HttpClient{
		config: config,
	}
}

func (h *HttpClient) GetAccrual(ctx context.Context, orderId int64) (response.OrderAccrual, error) {
	dto := response.OrderAccrual{}
	url := fmt.Sprintf("%v/api/orders/%v", h.config.AccrualUrl, orderId)
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return dto, err
	}

	res, err := client.Do(req)
	if err == nil {
		defer func(Body io.ReadCloser) {
			_, _ = io.Copy(io.Discard, res.Body)
			_ = Body.Close()
		}(res.Body)
	} else {
		return dto, err
	}

	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return dto, ErrNoOKResponse
	}

	err = json.Unmarshal(body, &dto)
	return dto, err
}
