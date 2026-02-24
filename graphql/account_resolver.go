package main

import (
	"context"
	"time"
)

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, obj *Account)([]*Order, error)  {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	orderList, err := r.server.orderClient.GetOrdersForAccount(ctx, obj.ID)
	if err != nil {
		return nil, err
	}
	orders := make([]*Order, len(orderList))

	for i, o := range orderList {

		products := make([]*OrderedProduct, len(o.Products))
		for j, p := range o.Products {
			products[j] = &OrderedProduct{
				ID: p.ID,
				Quantity: int(p.Quantity),
			}
		}

		orders[i] = &Order{
			ID: o.ID,
			TotalPrice: o.TotalPrice,
			CreatedAt: o.CreatedAt,
			Products: products,
		}
	}
	return orders, nil
}
