package grpc

import (
	"context"

	prodEntity "github.com/Mitra-Apps/be-store-service/domain/product/entity"
	prodPb "github.com/Mitra-Apps/be-store-service/domain/proto/product"
	pb "github.com/Mitra-Apps/be-store-service/domain/proto/store"
	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/Mitra-Apps/be-store-service/handler/grpc/middleware"
	"github.com/Mitra-Apps/be-store-service/service"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

type GrpcRoute struct {
	service service.Service
	pb.UnimplementedStoreServiceServer
}

func New(service service.Service) pb.StoreServiceServer {
	return &GrpcRoute{
		service: service,
	}
}

func (s *GrpcRoute) CreateStore(ctx context.Context, req *pb.CreateStoreRequest) (*pb.CreateStoreResponse, error) {
	if err := req.ValidateAll(); err != nil {
		st := status.New(codes.InvalidArgument, "Invalid argument").Proto()
		errsVal, ok := err.(pb.CreateStoreRequestMultiError)
		if ok {
			for _, err := range errsVal {
				errVal, ok := err.(pb.CreateStoreRequestValidationError)
				if ok {
					st.Details = append(st.Details, &anypb.Any{
						Value: []byte(errVal.Error()),
					})
				}
			}
		}

		return nil, err
	}

	store := &entity.Store{}
	if err := store.FromProto(req.Store); err != nil {
		return nil, err
	}

	store, err := s.service.CreateStore(ctx, store)
	if err != nil {
		return nil, err
	}

	return &pb.CreateStoreResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    store.ToProto(),
	}, nil
}

func (s *GrpcRoute) GetStore(ctx context.Context, req *pb.GetStoreRequest) (*pb.GetStoreResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	store, err := s.service.GetStore(ctx, req.StoreId)
	if err != nil {
		return nil, err
	}

	return &pb.GetStoreResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    store.ToProto(),
	}, nil
}

func (s *GrpcRoute) UpdateStore(ctx context.Context, req *pb.UpdateStoreRequest) (*pb.UpdateStoreResponse, error) {
	return nil, nil
}

func (s *GrpcRoute) DeleteStore(ctx context.Context, req *pb.DeleteStoreRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *GrpcRoute) ListStores(ctx context.Context, req *pb.ListStoresRequest) (*pb.ListStoresResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	stores, err := s.service.ListStores(ctx)
	if err != nil {
		return nil, err
	}

	result := &pb.ListStoresResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}

	for _, store := range stores {
		result.Data = append(result.Data, store.ToProto())
	}

	return result, nil
}

func (s *GrpcRoute) OpenCloseStore(ctx context.Context, req *pb.OpenCloseStoreRequest) (*pb.OpenCloseStoreResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Error when getting claims from jwt token")
	}

	err = s.service.OpenCloseStore(ctx, claims.UserID, claims.RoleNames, req.StoreId, req.IsActive)
	if err != nil {
		return nil, err
	}
	return &pb.OpenCloseStoreResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}, nil
}

func (s *GrpcRoute) CreateProducts(ctx context.Context, req *prodPb.CreateProductsRequest) (*prodPb.CreateProductsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	productList := []prodEntity.Product{}
	for _, p := range req.ProductList {
		pr := prodEntity.Product{}
		pr.FromProto(p)
		productList = append(productList, pr)
	}

	err := s.service.CreateProducts(ctx, productList)
	if err != nil {
		return nil, err
	}

	return &prodPb.CreateProductsResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}, nil
}
