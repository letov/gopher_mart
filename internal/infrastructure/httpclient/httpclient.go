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
	client *http.Client
	config *config.Config
}

func (h *HttpClient) GetAccrual(ctx context.Context, orderID string) (response.OrderAccrual, error) {
	url := fmt.Sprintf("%v/api/orders/%v", h.config.AccrualUrl, orderID)
	data, err := h.request(ctx, http.MethodGet, url, nil)
	if err != nil {
		return response.OrderAccrual{}, err
	}

	dto := response.OrderAccrual{}
	err = json.Unmarshal(data, &dto)
	return dto, err
}

func (h *HttpClient) request(
	ctx context.Context,
	method string,
	url string,
	body io.Reader,
) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return []byte{}, err
	}

	res, err := h.client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer func(Body io.ReadCloser) {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = Body.Close()
	}(res.Body)

	rb, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return []byte{}, ErrNoOKResponse
	}

	return rb, err
}

func NewHttpClient(
	config *config.Config,
) *HttpClient {
	client := &http.Client{}
	return &HttpClient{
		client,
		config,
	}
}
