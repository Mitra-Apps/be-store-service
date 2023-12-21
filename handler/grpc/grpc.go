package grpc

import (
	"context"

	pb "github.com/Mitra-Apps/be-store-service/domain/proto/store"
	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/Mitra-Apps/be-store-service/service"
	"github.com/golang/protobuf/ptypes/empty"
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

func (s *GrpcRoute) CreateStore(ctx context.Context, req *pb.CreateStoreRequest) (*pb.Store, error) {
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

	return store.ToProto(), nil
}

func (s *GrpcRoute) GetStore(ctx context.Context, req *pb.GetStoreRequest) (*pb.Store, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	store, err := s.service.GetStore(ctx, req.StoreId)
	if err != nil {
		return nil, err
	}

	return store.ToProto(), nil
}

func (s *GrpcRoute) UpdateStore(ctx context.Context, req *pb.UpdateStoreRequest) (*pb.Store, error) {
	return nil, nil
}

func (s *GrpcRoute) DeleteStore(ctx context.Context, req *pb.DeleteStoreRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *GrpcRoute) ListStores(ctx context.Context, req *pb.ListStoresRequest) (*pb.ListStoresResponse, error) {
	return nil, nil
}
