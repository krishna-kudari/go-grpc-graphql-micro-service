package order

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/krishna-kudari/go-grpc-graphql-micro-service/account"
	"github.com/krishna-kudari/go-grpc-graphql-micro-service/catalog"
	pb "github.com/krishna-kudari/go-grpc-graphql-micro-service/order/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	accountClientTimeout = 5 * time.Second
	catalogClientTimeout = 5 * time.Second
)

type grpcServer struct {
	pb.UnimplementedOrderServiceServer
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
	logger        *slog.Logger
}

// ListenGRPC starts a gRPC server listening on the specified port.
// It returns an error if initialization fails, but does not block.
// The caller is responsible for graceful shutdown.
func ListenGRPC(ctx context.Context, service Service, accountURL, catalogURL string, port int, logger *slog.Logger) error {
	if logger == nil {
		logger = slog.Default()
	}

	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return fmt.Errorf("creating account client: %w", err)
	}

	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return fmt.Errorf("creating catalog client: %w", err)
	}

	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return fmt.Errorf("listening on %s: %w", addr, err)
	}

	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, &grpcServer{
		service:       service,
		accountClient: accountClient,
		catalogClient: catalogClient,
		logger:        logger,
	})
	reflection.Register(server)

	logger.Info("gRPC server listening", slog.String("address", addr))
	return server.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	if r == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if r.AccountId == "" {
		return nil, status.Error(codes.InvalidArgument, "account_id is required")
	}
	if len(r.Products) == 0 {
		return nil, status.Error(codes.InvalidArgument, "products cannot be empty")
	}

	accountCtx, cancel := context.WithTimeout(ctx, accountClientTimeout)
	defer cancel()

	_, err := s.accountClient.GetAccountInfo(accountCtx, r.AccountId)
	if err != nil {
		s.logger.Error("failed to get account info",
			slog.String("account_id", r.AccountId),
			slog.String("error", err.Error()))
		return nil, status.Errorf(codes.NotFound, "account not found: %v", err)
	}

	productIDs := make([]string, len(r.Products))
	for i, p := range r.Products {
		if p == nil {
			return nil, status.Error(codes.InvalidArgument, "product cannot be nil")
		}
		productIDs[i] = p.ProductId
	}

	catalogCtx, cancel := context.WithTimeout(ctx, catalogClientTimeout)
	defer cancel()

	orderedProducts, err := s.catalogClient.GetProductsByIDs(catalogCtx, productIDs)
	if err != nil {
		s.logger.Error("failed to get products",
			slog.Any("product_ids", productIDs),
			slog.String("error", err.Error()))
		return nil, status.Errorf(codes.NotFound, "products not found: %v", err)
	}

	products := make([]OrderedProduct, 0, len(r.Products))
	productMap := make(map[string]uint32)
	for _, rp := range r.Products {
		if rp != nil {
			productMap[rp.ProductId] = rp.Quantity
		}
	}

	for _, p := range orderedProducts {
		if quantity, exists := productMap[p.ID]; exists && quantity > 0 {
			products = append(products, OrderedProduct{
				ID:        p.ID,
				ProductID: p.ID,
				Quantity:  quantity,
			})
		}
	}

	if len(products) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no valid products with quantity > 0")
	}

	order, err := s.service.PostOrder(ctx, r.AccountId, products)
	if err != nil {
		s.logger.Error("failed to create order",
			slog.String("account_id", r.AccountId),
			slog.String("error", err.Error()))
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	return &pb.PostOrderResponse{
		Order: &pb.Order{
			Id: order.ID,
		},
	}, nil
}

func (s *grpcServer) GetOrdersForAccount(ctx context.Context, r *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
	if r == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if r.AccountId == "" {
		return nil, status.Error(codes.InvalidArgument, "account_id is required")
	}

	orders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)
	if err != nil {
		s.logger.Error("failed to get orders for account",
			slog.String("account_id", r.AccountId),
			slog.String("error", err.Error()))
		return nil, status.Errorf(codes.Internal, "failed to get orders: %v", err)
	}

	pbOrders := make([]*pb.Order, len(orders))
	for i, o := range orders {
		pbOrder, err := toProtoOrder(o)
		if err != nil {
			s.logger.Error("failed to convert order to proto",
				slog.String("order_id", o.ID),
				slog.String("error", err.Error()))
			return nil, status.Errorf(codes.Internal, "failed to convert order: %v", err)
		}
		pbOrders[i] = pbOrder
	}

	return &pb.GetOrdersForAccountResponse{
		Orders: pbOrders,
	}, nil
}
