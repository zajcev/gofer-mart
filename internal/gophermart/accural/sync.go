package accural

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zajcev/gofer-mart/internal/gophermart/config"
	"github.com/zajcev/gofer-mart/internal/gophermart/database"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"io"
	"log"
	"net/http"
)

func AccuralIntegration() {
	list, err := database.GetActiveOrders(context.Background())
	if err != nil {
		log.Printf("Error while GetActiveOrders: %v", err)
		return
	}
	for _, v := range list {
		order, err := sendToAccuralSystem(&v)
		if err != nil {
			log.Printf("Error while sendToAccuralSystem: %v", err)
			continue
		}
		if order.Status != "" {
			updateOrderStatus(context.Background(), order)
			if order.Status == "PROCESSED" {
				updateOrderAccural(context.Background(), order)
			}
		}
	}
}

func sendToAccuralSystem(o *model.Order) (*model.Order, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", config.GetAccSystemAddr()+"/api/orders/"+o.ID, nil)
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

func updateOrderStatus(ctx context.Context, o *model.Order) {
	if o.Status == "INVALID" {
		log.Printf("Order %s is invalid", o.ID)
		database.UpdateOrderStatus(ctx, o)
	}
	if o.Status == "REGISTERED" || o.Status == "PROCESSING" {
		log.Printf("Order %v in status %v", o.ID, o.Status)
		database.UpdateOrderStatus(ctx, o)
	}
	if o.Status == "PROCESSED" {
		database.UpdateOrderStatus(ctx, o)
	}
}

func updateOrderAccural(ctx context.Context, o *model.Order) {
	database.UpdateOrderAccural(ctx, o)
	err := database.SetCurrent(ctx, o)
	if err != nil {
		return
	}
}
