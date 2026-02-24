package main

import (
	"context"
	"fmt"
	"time"

	"github.com/krishna-kudari/go-grpc-graphql-micro-service/order"
)

type mutationResolver struct{
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, in AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	a, err := r.server.accountClient.CreateAccount(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return &Account{
		ID: a.ID,
		Name: a.Name,
	}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	p, err := r.server.catalogClient.PostProduct(ctx, in.Name, in.Description, float32(in.Price))
	if err != nil {
		return nil, err
	}

	return &Product{
		ID: p.ID,
		Description: p.Description,
		Name: p.Name,
		Price: p.Price,
	}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	products :=	make([]order.ProductPayload, len(in.Products))
	for i, p := range in.Products {
		if p.Quantity <= 0 {
			return nil, fmt.Errorf("invalid param quantity : %d",p.Quantity)
		}
		products[i] = order.ProductPayload{
			ID: p.ID,
			Quantity: uint32(p.Quantity),
		}
	}
	o, err := r.server.orderClient.CreateOrder(ctx, in.AccountID, products)
	if err != nil {
		return nil, err
	}

	orderedProducts := make([]*OrderedProduct, len(o.Products))
	for i, op := range o.Products {
		orderedProducts[i] = &OrderedProduct{
			ID: op.ID,
			Quantity: int(op.Quantity),
		}
	}
	return &Order{
		ID: o.ID,
		CreatedAt: o.CreatedAt,
		TotalPrice: o.TotalPrice,
		Products: orderedProducts,
	}, nil
}
