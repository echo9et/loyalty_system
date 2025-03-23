package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	config "gophermart.ru/internal"
)

type AccrualSystem struct {
	addr   string
	client *http.Client
}

type OrderResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

func (a *AccrualSystem) Do(req *http.Request) (*http.Response, error) {
	return a.client.Do(req)
}

func (a *AccrualSystem) GetOrderInfo(orderNumber string) (*OrderResponse, int, error) {
	url := fmt.Sprintf("http://%s/api/orders/%s", a.addr, orderNumber)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := a.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, err
	}

	var orderResp OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return nil, 0, err
	}

	return &orderResp, resp.StatusCode, nil
}

func NewAccrualSystem() *AccrualSystem {
	return &AccrualSystem{
		addr:   config.Get().AddrAccraulSystem,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (as *AccrualSystem) CheckOrder(orderID string) error {
	return nil
}
