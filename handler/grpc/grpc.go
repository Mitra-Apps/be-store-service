package grpc

import (
	"context"

	pb "github.com/Mitra-Apps/be-store-service/domain/proto/store"
	"github.com/Mitra-Apps/be-store-service/service"
)

type GrpcRoute struct {
	service service.ServiceInterface
	pb.UnimplementedStoreServiceServer
}

func New(service service.ServiceInterface) pb.StoreServiceServer {
	return &GrpcRoute{
		service: service,
	}
}

func (g *GrpcRoute) GetStores(ctx context.Context, req *pb.GetStoresRequest) (*pb.GetStoresResponse, error) {
	stores, err := g.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	protoStores := []*pb.Store{}
	for _, store := range stores {
		protoStores = append(protoStores, store.ToProto())
	}

	return &pb.GetStoresResponse{
		Stores: protoStores,
	}, nil
}
