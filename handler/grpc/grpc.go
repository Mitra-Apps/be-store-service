package grpc

import (
	"context"

	prodEntity "github.com/Mitra-Apps/be-store-service/domain/product/entity"
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

func (s *GrpcRoute) UpsertProducts(ctx context.Context, req *pb.UpsertProductsRequest) (*pb.UpsertProductsResponse, error) {
	if err := validateProduct(req.ProductList); err != nil {
		return nil, err
	}

	productList := []prodEntity.Product{}
	for _, p := range req.ProductList {
		pr := prodEntity.Product{}
		if err := pr.FromProto(p); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		productList = append(productList, pr)
	}

	err := s.service.UpsertProducts(ctx, productList)
	if err != nil {
		return nil, err
	}

	return &pb.UpsertProductsResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}, nil
}

func validateProduct(products []*pb.Product) error {
	for _, p := range products {
		if p.StoreId == "" {
			return status.Errorf(codes.InvalidArgument, "Store id is required")
		}
		if p.Name == "" {
			return status.Errorf(codes.InvalidArgument, "Name is required")
		}
		if p.Price <= 0 {
			return status.Errorf(codes.InvalidArgument, "Price is required")
		}
		if p.UomId == "" {
			return status.Errorf(codes.InvalidArgument, "unit of measure is required")
		}
		if p.ProductTypeId == "" {
			return status.Errorf(codes.InvalidArgument, "product type id is required")
		}
	}
	return nil
}

func (g *GrpcRoute) UpsertUnitOfMeasure(ctx context.Context, req *pb.UpsertUnitOfMeasureRequest) (*pb.UpsertUnitOfMeasureResponse, error) {
	if req.Uom.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}
	if req.Uom.Symbol == "" {
		return nil, status.Errorf(codes.InvalidArgument, "symbol is required")
	}

	if err := g.service.UpsertUnitOfMeasure(ctx, prodEntity.UnitOfMeasure{
		Name:     req.Uom.Name,
		Symbol:   req.Uom.Symbol,
		IsActive: req.Uom.IsActive,
	}); err != nil {
		return nil, err
	}

	return &pb.UpsertUnitOfMeasureResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}, nil
}

func (g *GrpcRoute) UpsertProductCategory(ctx context.Context, req *pb.UpsertProductCategoryRequest) (*pb.UpsertProductCategoryResponse, error) {
	if req.ProductCategory.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	return &pb.UpsertProductCategoryResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}, nil
}

func (g *GrpcRoute) UpsertProductType(ctx context.Context, req *pb.UpsertProductTypeRequest) (*pb.UpsertProductTypeResponse, error) {
	if req.ProductType.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	if req.ProductType.ProductCategoryId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "product category id is required")
	}

	return &pb.UpsertProductTypeResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}, nil
}
