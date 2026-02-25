package catalog

import (
	"context"
	"fmt"
	"net"

	"github.com/krishna-kudari/go-grpc-graphql-micro-service/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedCatalogServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterCatalogServiceServer(server, &grpcServer{
		service:                           s,
		UnimplementedCatalogServiceServer: pb.UnimplementedCatalogServiceServer{},
	})
	reflection.Register(server)
	return server.Serve(listener)
}

func (s *grpcServer) PostProduct(ctx context.Context,r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	res, err := s.service.PostProduct(ctx, r.Name, r.Description, float64(r.Price))
	if err != nil {
		return nil, err
	}
	return &pb.PostProductResponse{
		Id: res.ID,
		Name: res.Name,
		Description: res.Description,
		Price: float32(res.Price),
	}, nil
}

func (s *grpcServer) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	res, err := s.service.GetProduct(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{
		Id: res.ID,
		Name: res.Name,
		Description: res.Description,
		Price: float32(res.Price),
	}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	res, err := s.service.GetProducts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}
	var products []*pb.Product
	for _, product := range res {
		products = append(products, &pb.Product{
			Id: product.ID,
			Name: product.Name,
			Description: product.Description,
			Price: float32(product.Price),
		})
	}
	return &pb.GetProductsResponse{
		Products: products,
	}, nil
}
func (s *grpcServer) GetProductsByIDs(ctx context.Context, r *pb.GetProductsByIDsRequest) (*pb.GetProductsByIDsResponse, error) {
	res, err := s.service.GetProductsByIDs(ctx, r.Ids)
	if err != nil {
		return nil, err
	}
	var products []*pb.Product
	for _, product := range res {
		products = append(products, &pb.Product{
			Id: product.ID,
			Name: product.Name,
			Description: product.Description,
			Price: float32(product.Price),
		})
	}
	return &pb.GetProductsByIDsResponse{
		Products: products,
	}, nil
}
func (s *grpcServer) SearchProducts(ctx context.Context, r *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	res, err := s.service.SearchProducts(ctx, r.Query, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}
	var products []*pb.Product
	for _, product := range res {
		products = append(products, &pb.Product{
			Id: product.ID,
			Name: product.Name,
			Description: product.Description,
			Price: float32(product.Price),
		})
	}
	return &pb.SearchProductsResponse{
		Products: products,
	}, nil
}
