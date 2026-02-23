package catalog

import (
	"context"
	"log"

	"github.com/krishna-kudari/go-grpc-graphql-micro-service/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	serviceClient := pb.NewCatalogServiceClient(conn)
	return &Client{
		conn:    conn,
		service: serviceClient,
	}, nil
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) PostProduct(ctx context.Context, name string, description string, price float32) (*Product, error) {
	r, err := c.service.PostProduct(
		ctx,
		&pb.PostProductRequest{
			Name:        name,
			Description: description,
			Price:       price,
		},
	)
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		Price:       float64(r.Price),
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	r, err := c.service.GetProduct(ctx, &pb.GetProductRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	product := &Product{
		ID:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		Price:       float64(r.Price),
	}
	return product, nil
}

func (c *Client) GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	res, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{
		Skip: skip,
		Take: take,
	})
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for _, p := range res.Products {
		products = append(products, Product{
			Name:        p.Name,
			ID:          p.Id,
			Description: p.Description,
			Price:       float64(p.Price),
		})
	}
	return products, nil
}

func (c *Client) GetProductsByIDs(ctx context.Context, ids []string) ([]Product, error) {
	res, err := c.service.GetProductsByIDs(ctx, &pb.GetProductsByIDsRequest{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for _, p := range res.Products {
		products = append(products, Product{
			Name:        p.Name,
			ID:          p.Id,
			Description: p.Description,
			Price:       float64(p.Price),
		})
	}
	return products, nil
}

func (c *Client) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	res, err := c.service.SearchProducts(ctx, &pb.SearchProductsRequest{
		Query: query,
		Skip:  skip,
		Take:  take,
	})
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for _, p := range res.Products {
		products = append(products, Product{
			Name:        p.Name,
			ID:          p.Id,
			Description: p.Description,
			Price:       float64(p.Price),
		})
	}
	return products, nil
}
