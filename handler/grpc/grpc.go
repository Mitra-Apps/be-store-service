package grpc

import (
	"context"

	pb "github.com/Mitra-Apps/be-store-service/domain/proto/store"
	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/Mitra-Apps/be-store-service/handler/grpc/middleware"
	"github.com/Mitra-Apps/be-store-service/service"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	store := &entity.Store{}
	if err := store.FromProto(req.Store); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	data, err := s.service.UpdateStore(ctx, req.StoreId, store)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateStoreResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data.ToProto(),
	}, nil
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
