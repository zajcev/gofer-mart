package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"io"
	"log"
	"net/http"
	"time"
)

type AccrualStorage interface {
	UpdateOrderStatus(ctx context.Context, o *model.Order) int
	UpdateOrderAccrual(ctx context.Context, o *model.Order) int
	GetActiveOrders(ctx context.Context) ([]model.Order, error)
	SetCurrent(ctx context.Context, order *model.Order) error
}
type Accrual struct {
	db AccrualStorage
}

func NewAccrual(db AccrualStorage) *Accrual {
	return &Accrual{db: db}
}

func (a *Accrual) AccrualIntegration(ctx context.Context, url string) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			list, err := a.db.GetActiveOrders(ctx)
			if err != nil {
				return fmt.Errorf("failed to get active orders: %w", err)
			}
			for _, v := range list {
				order, err := sendToAccrualSystem(&v, url)
				if err != nil {
					log.Printf("error sending order %s to accrual: %v", v.ID, err)
					continue
				}
				if order.Status != "" {
					updateOrderStatus(ctx, order, a)
					if order.Status == "PROCESSED" {
						updateOrderAccrual(ctx, order, a)
						log.Printf("error updating order %s accrual: %v", order.ID, err)
					}
				}
			}
		}
	}
}

func sendToAccrualSystem(o *model.Order, url string) (*model.Order, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url+"/api/orders/"+o.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request failed after retries: %v", err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("accrual system returned status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	if len(body) == 0 {
		return nil, errors.New("empty response body")
	}

	err = json.Unmarshal(body, &o)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}
	return o, nil
}

func updateOrderStatus(ctx context.Context, o *model.Order, a *Accrual) {
	a.db.UpdateOrderStatus(ctx, o)
}

func updateOrderAccrual(ctx context.Context, o *model.Order, a *Accrual) {
	a.db.UpdateOrderAccrual(ctx, o)
	err := a.db.SetCurrent(ctx, o)
	if err != nil {
		return
	}
}
