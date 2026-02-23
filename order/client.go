package order

import (
	"context"
	"fmt"

	pb "github.com/krishna-kudari/go-grpc-graphql-micro-service/order/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


type Client struct {
	connection *grpc.ClientConn
	service pb.OrderServiceClient
}


func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	serviceClient := pb.NewOrderServiceClient(conn)
	return &Client{
		connection: conn,
		service: serviceClient,
	}, nil
}

type ProductPayload struct {
	ID string `json:"id"`
	Quantity uint32 `json:"quantity"`
}

func (c *Client) CreateOrder(ctx context.Context, accountID string, products []ProductPayload) (*Order, error) {
	if accountID == "" {
		return nil, fmt.Errorf("order.CreateOrder: accountID must not be empty")
	}
	if len(products) <= 0 {
		return nil, fmt.Errorf("order.CreateOrder: products must not be empty")
	}
	req := &pb.PostOrderRequest{
		AccountId: accountID,
		Products: toProtoProducts(products),
	}

	res, err := c.service.PostOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("order.CreateOrder: post order rpc: %w", err)
	}
	if res == nil {
		return nil, fmt.Errorf("order.CreateOrder: post order rpc server returned nil order")
	}
	order, err := toDomainOrder(res.Order)
	if err != nil {
		return nil, fmt.Errorf("order.CreateOrder: converting response: %w", err)
	}
	return order,nil
}

func (c *Client) GetOrdersForAccount(ctx context.Context, accountID string) ([]*Order, error) {
	if accountID == "" {
		return nil, fmt.Errorf("order.GetOrdersForAccount: accountID must not be empty")
	}
	req := &pb.GetOrdersForAccountRequest{
		AccountId: accountID,
	}
	res, err := c.service.GetOrdersForAccount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("order.GetOrdersForAccount: GetOrdersForAccount rpc err %w",err)
	}
	if res.Orders == nil {
		return nil, fmt.Errorf("order.GetOrdersForAccount: GetOrdersForAccount rpc server returned nil orders")
	}

	orders, err := toDomainOrders(res.Orders)
	if err != nil {
		return nil, fmt.Errorf("order.GetOrdersForAccount: converting response: %w", err)
	}
	return orders, nil
}
