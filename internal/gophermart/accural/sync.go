package accural

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zajcev/gofer-mart/internal/gophermart/database"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"io"
	"log"
	"net/http"
	"time"
)

type Accrual struct {
	db *database.DBService
}

func NewAccrual(db *database.DBService) *Accrual {
	return &Accrual{db: db}
}

func (a *Accrual) AccrualIntegration(ctx context.Context, url string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
			list, err := a.db.GetActiveOrders(context.Background())
			if err != nil {
				return err
			}
			for _, v := range list {
				order, err := sendToAccrualSystem(&v, url)
				if err != nil {
					continue
				}
				if order.Status != "" {
					updateOrderStatus(context.Background(), order, a)
					if order.Status == "PROCESSED" {
						updateOrderAccrual(context.Background(), order, a)
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
	println("accrual system successfully sent")
	return o, nil
}

func updateOrderStatus(ctx context.Context, o *model.Order, a *Accrual) {
	a.db.UpdateOrderStatus(ctx, o)
}

func updateOrderAccrual(ctx context.Context, o *model.Order, a *Accrual) {
	a.db.UpdateOrderAccural(ctx, o)
	err := a.db.SetCurrent(ctx, o)
	if err != nil {
		return
	}
}
