package order

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
)

// TimeSource provides the current time. Used for testing.
type TimeSource interface {
	Now() time.Time
}

type realTimeSource struct{}

func (realTimeSource) Now() time.Time {
	return time.Now()
}

type OrderedProduct struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Quantity  uint32 `json:"quantity"`
	OrderID   string `json:"order_id"`
}

type Order struct {
	ID         string           `json:"id"`
	CreatedAt  time.Time        `json:"created_at"`
	AccountID  string           `json:"account_id"`
	TotalPrice float64          `json:"total_price"`
	Products   []OrderedProduct `json:"products"`
}

type Service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type orderService struct {
	repository Repository
	timeSource TimeSource
}

func (o *orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	if accountID == "" {
		return nil, fmt.Errorf("get orders for account: account_id cannot be empty")
	}

	orders, err := o.repository.GetOrdersForAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get orders for account %q: %w", accountID, err)
	}
	return orders, nil
}

func (o *orderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	if accountID == "" {
		return nil, fmt.Errorf("post order: account_id cannot be empty")
	}
	if len(products) == 0 {
		return nil, fmt.Errorf("post order: products cannot be empty")
	}

	for i, p := range products {
		if p.Quantity == 0 {
			return nil, fmt.Errorf("post order: product at index %d has zero quantity", i)
		}
		if p.ProductID == "" {
			return nil, fmt.Errorf("post order: product at index %d has empty product_id", i)
		}
	}

	order := Order{
		ID:         ksuid.New().String(),
		AccountID:  accountID,
		TotalPrice: 0, // TODO: calculate total price from products
		Products:   products,
		CreatedAt:  o.timeSource.Now(),
	}

	for i := range order.Products {
		order.Products[i].OrderID = order.ID
	}

	if err := o.repository.PutOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("post order: saving order: %w", err)
	}
	return &order, nil
}

func NewOrderService(r Repository) (Service, error) {
	if r == nil {
		return nil, fmt.Errorf("new order service: repository cannot be nil")
	}
	return &orderService{
		repository: r,
		timeSource: realTimeSource{},
	}, nil
}

func NewOrderServiceWithTimeSource(r Repository, ts TimeSource) (Service, error) {
	if r == nil {
		return nil, fmt.Errorf("new order service: repository cannot be nil")
	}
	if ts == nil {
		return nil, fmt.Errorf("new order service: time source cannot be nil")
	}
	return &orderService{
		repository: r,
		timeSource: ts,
	}, nil
}
