package order

import (
	"fmt"

	pb "github.com/krishna-kudari/go-grpc-graphql-micro-service/order/pb/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoProducts(products []ProductPayload) []*pb.OrderedProduct {
	out := make([]*pb.OrderedProduct, len(products))
	for i, p := range products {
		out[i] = &pb.OrderedProduct{
			ProductId: p.ID,
			Quantity:  p.Quantity,
		}
	}
	return out
}

func toDomainOrder(o *pb.Order) (*Order, error) {
	if o == nil {
		return nil, fmt.Errorf("toDomainOrder: received nil proto order")
	}
	if o.CreatedAt == nil {
		return nil, fmt.Errorf("toDomainOrder: order %q has nil CreatedAt", o.Id)
	}

	products := make([]OrderedProduct, len(o.Products))
	for i, p := range o.Products {
		products[i] = OrderedProduct{
			ID:        p.Id,
			ProductID: p.ProductId,
			OrderID:   p.OrderId,
			Quantity:  p.Quantity,
		}
	}
	return &Order{
		ID:         o.Id,
		AccountID:  o.AccountId,
		TotalPrice: o.TotalPrice,
		CreatedAt:  o.CreatedAt.AsTime(),
		Products:   products,
	}, nil
}

func toDomainOrders(orderes []*pb.Order) ([]*Order, error) {
	out := make([]*Order, len(orderes))
	for i, o := range orderes {
		order, err := toDomainOrder(o)
		if err != nil {
			return nil, err
		}
		out[i] = order
	}
	return out, nil
}

func toProtoOrder(o Order) (*pb.Order, error) {
	createdAt := timestamppb.New(o.CreatedAt)
	if err := createdAt.CheckValid(); err != nil {
		return nil, fmt.Errorf("converting created_at to timestamp: %w", err)
	}

	products := make([]*pb.OrderedProduct, len(o.Products))
	for i, p := range o.Products {
		products[i] = &pb.OrderedProduct{
			Id:        p.ID,
			OrderId:   p.OrderID,
			ProductId: p.ProductID,
			Quantity:  p.Quantity,
		}
	}

	return &pb.Order{
		Id:         o.ID,
		AccountId:  o.AccountID,
		TotalPrice: o.TotalPrice,
		CreatedAt:  createdAt,
		Products:   products,
	}, nil
}
