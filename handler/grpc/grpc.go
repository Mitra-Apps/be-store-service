package grpc

import (
	"context"
	"fmt"
	"strings"

	prodEntity "github.com/Mitra-Apps/be-store-service/domain/product/entity"
	errPb "github.com/Mitra-Apps/be-store-service/domain/proto"
	pb "github.com/Mitra-Apps/be-store-service/domain/proto/store"
	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/Mitra-Apps/be-store-service/handler/grpc/middleware"
	"github.com/Mitra-Apps/be-store-service/service"
	util "github.com/Mitra-Apps/be-utility-service/service"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
		st := status.New(codes.InvalidArgument, "Error when validating request")
		if multiErrs, ok := err.(pb.CreateStoreRequestMultiError); ok {
			for _, multiErr := range multiErrs {
				if validationErr, ok := multiErr.(pb.CreateStoreRequestValidationError); ok {
					// print type of validationErr error
					fmt.Printf("%T\n", validationErr.Cause())
					detail := errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequest_FieldViolation{
							{
								Field:       validationErr.Field(),
								Description: validationErr.Cause().Error(),
							},
						},
					}
					st, _ = st.WithDetails(&detail)
				}
			}
		}

		return nil, st.Err()
	}

	storeHoursErr := validateStoreHours(req.Store.Hours)
	if storeHoursErr != nil {
		return nil, storeHoursErr
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

	if store == nil {
		return &pb.GetStoreResponse{
			Code:    int32(codes.NotFound),
			Message: "Store not found",
			Data:    nil,
		}, nil
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

	storeHoursErr := validateStoreHours(req.Store.Hours)
	if storeHoursErr != nil {
		return nil, storeHoursErr
	}

	store := &entity.Store{}
	if err := store.FromProto(req.Store); err != nil {
		return nil, err
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
	for _, id := range req.GetIds() {
		_, err := uuid.Parse(id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Error when parsing store id to uuid")
		}
	}

	return &emptypb.Empty{}, s.service.DeleteStores(ctx, req.GetIds())
}

func (s *GrpcRoute) ListStores(ctx context.Context, req *pb.ListStoresRequest) (*pb.ListStoresResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	stores, err := s.service.ListStores(ctx, 1, 20)
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

func (s *GrpcRoute) GetStoreByUserID(ctx context.Context, req *pb.GetStoreByUserIDRequest) (*pb.GetStoreByUserIDResponse, error) {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Error when getting claims from jwt token")
	}
	store, err := s.service.GetStoreByUserID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if store == nil {
		return &pb.GetStoreByUserIDResponse{
			Code:    int32(codes.OK),
			Message: "Store not found",
		}, nil
	}

	return &pb.GetStoreByUserIDResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    store.ToProto(),
	}, nil
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

func (s *GrpcRoute) InsertProducts(ctx context.Context, req *pb.InsertProductsRequest) (*pb.GenericResponse, error) {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, util.NewError(codes.Unauthenticated, errPb.StoreErrorCode_ERROR_WHEN_GETTING_CLAIMS_FROM_JWT_TOKEN.String(), "Error saat mendapatkan claims dari jwt token")
	}

	productList := []*prodEntity.Product{}
	for _, p := range req.ProductList {
		pr := prodEntity.Product{}
		if err := pr.FromProto(p, &req.StoreId); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		productList = append(productList, &pr)
	}

	if err := validateProduct(productList...); err != nil {
		return nil, err
	}

	storeIdUuid, err := uuid.Parse(req.StoreId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = s.service.UpsertProducts(ctx, claims.UserID, claims.RoleNames, storeIdUuid, false, productList...)
	if err != nil {
		return nil, err
	}

	return &pb.GenericResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}, nil
}

func (s *GrpcRoute) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.GenericResponse, error) {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, util.NewError(codes.Unauthenticated, errPb.StoreErrorCode_ERROR_WHEN_GETTING_CLAIMS_FROM_JWT_TOKEN.String(), "Error saat mendapatkan claims dari jwt token")
	}

	product := &prodEntity.Product{}
	req.Product.Id = req.ProductId
	if err := product.FromProto(req.Product, nil); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	if err := validateProduct(product); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = s.service.UpsertProducts(ctx, claims.UserID, claims.RoleNames, product.StoreID, true, product)
	if err != nil {
		return nil, err
	}

	return &pb.GenericResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}, nil
}

func (s *GrpcRoute) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*empty.Empty, error) {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, util.NewError(codes.Unauthenticated, errPb.StoreErrorCode_ERROR_WHEN_GETTING_CLAIMS_FROM_JWT_TOKEN.String(), "Error saat mendapatkan claims dari jwt token")
	}

	id, err := uuid.Parse(req.GetProductId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error when parsing store id to uuid")
	}

	return &emptypb.Empty{}, s.service.DeleteProductById(ctx, claims.UserID, id)
}

func validateProduct(products ...*prodEntity.Product) error {
	for _, p := range products {
		if p.Name == "" {
			return util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_NAME_IS_REQUIRED.String(), "Nama tidak boleh kosong")
		}
		if p.Price <= 0 {
			return util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_PRICE_IS_REQUIRED.String(), "Harga tidak boleh kosong")
		}
		if p.Uom == "" {
			return util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_UOM_IS_REQUIRED.String(), "Satuan unit diperlukan")
		}
		if p.ProductTypeID == 0 {
			return util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_PRODUCT_TYPE_IS_REQUIRED.String(), "product type diperlukan")
		}
	}
	return nil
}

func (g *GrpcRoute) UpsertProductCategory(ctx context.Context, req *pb.UpsertProductCategoryRequest) (*pb.UpsertProductCategoryResponse, error) {
	//TODO: refactor the validation to another function to minimize duplicate code
	isUpdateCategory := false
	if req.GetId() > 0 {
		req.ProductCategory.Id = req.GetId()
		isUpdateCategory = true
	}

	if err := req.ValidateAll(); err != nil {
		st := status.New(codes.InvalidArgument, "Error when validating request")
		if multiErrs, ok := err.(pb.UpsertProductCategoryRequestMultiError); ok {
			for _, multiErr := range multiErrs {
				if validationErr, ok := multiErr.(pb.UpsertProductCategoryRequestValidationError); ok {
					detail := errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequest_FieldViolation{
							{
								Field:       validationErr.Field(),
								Description: validationErr.Cause().Error(),
							},
						},
					}
					st, _ = st.WithDetails(&detail)
				}
			}
		}

		return nil, st.Err()
	}

	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Error when getting claims from jwt token")
	}

	prodCat := new(prodEntity.ProductCategory)
	prodCat.FromProto(req.ProductCategory)
	prodCat.CreatedBy = claims.UserID

	if !isUpdateCategory {
		if err := g.service.UpsertProductCategory(ctx, prodCat); err != nil {
			return nil, err
		}
	} else {
		if err := g.service.UpdateProductCategory(ctx, prodCat); err != nil {
			return nil, err
		}
	}

	return &pb.UpsertProductCategoryResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}, nil
}

func (g *GrpcRoute) UpdateProductCategory(ctx context.Context, req *pb.UpsertProductCategoryRequest) (*pb.UpsertProductCategoryResponse, error) {
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	return g.UpsertProductCategory(ctx, req)
}

func (g *GrpcRoute) UpsertProductType(ctx context.Context, req *pb.UpsertProductTypeRequest) (*pb.UpsertProductTypeResponse, error) {

	prodType, err := g.validateUpdateProductTypeRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := g.service.UpsertProductType(ctx, prodType); err != nil {
		return nil, err
	}

	return &pb.UpsertProductTypeResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}, nil
}

func (g *GrpcRoute) GetProductById(ctx context.Context, req *pb.GetProductByIdRequest) (*pb.GetProductByIdResponse, error) {
	prodId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error when parsing product id to uuid")
	}
	prod, err := g.service.GetProductById(ctx, prodId)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductByIdResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    prod.ToProto(),
	}, nil
}

func (g *GrpcRoute) GetProductList(ctx context.Context, req *pb.GetProductListRequest) (*pb.GetProductListResponse, error) {
	if req.Page == 0 || req.Limit == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Page and Limit is required")
	}

	if strings.Trim(req.StoreId, " ") == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Store id is required")
	}

	storeId, err := uuid.Parse(req.StoreId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error when parsing store id to uuid")
	}

	var productTypeId *int64
	if req.ProductTypeId != 0 {
		productTypeId = &req.ProductTypeId
	}

	products, pagination, err := g.service.GetProductsByStoreId(ctx, req.Page, req.Limit, storeId, productTypeId, req.IsIncludeDeactivated)
	if err != nil {
		return nil, err
	}

	data := []*pb.Product{}
	for _, p := range products {
		data = append(data, p.ToProto())
	}

	paging := pb.Pagination{
		Page:        pagination.Page,
		Limit:       pagination.Limit,
		TotalRecord: pagination.TotalRecords,
		Records:     pagination.Records,
		TotalPage:   pagination.TotalPage,
	}

	return &pb.GetProductListResponse{
		Code:       int32(codes.OK),
		Message:    codes.OK.String(),
		Data:       data,
		Pagination: &paging,
	}, nil
}

func (g *GrpcRoute) GetProductCategories(ctx context.Context, req *pb.GetProductCategoriesRequest) (*pb.GetProductCategoriesResponse, error) {
	cat, uom, err := g.service.GetProductCategories(ctx, req.IsIncludeDeactivated)
	if err != nil {
		return nil, err
	}
	cats := []*pb.ProductCategory{}
	for _, u := range cat {
		cats = append(cats, u.ToProto())
	}
	return &pb.GetProductCategoriesResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &pb.GetProductCategoriesResponseItem{
			ProductCategory: cats,
			Uom:             uom,
		},
	}, nil
}

func (g *GrpcRoute) GetProductTypes(ctx context.Context, req *pb.GetProductTypesRequest) (*pb.GetProductTypesResponse, error) {
	prodType, err := g.service.GetProductTypes(ctx, req.ProductCategoryId, req.IsIncludeDeactivated)
	if err != nil {
		return nil, err
	}
	prodTypes := []*pb.ProductType{}
	for _, u := range prodType {
		prodTypes = append(prodTypes, u.ToProto())
	}
	return &pb.GetProductTypesResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    prodTypes,
	}, nil
}

func validateStoreHours(hours []*pb.StoreHour) error {
	DAYS := map[int32]string{
		0: "MONDAY",
		1: "TUESDAY",
		2: "WEDNESDAY",
		3: "THURSDAY",
		4: "FRIDAY",
		5: "SATURDAY",
		6: "SUNDAY",
	}

	errs := []string{}

	for _, hour := range hours {
		if !hour.IsOpen && hour.Is24Hours {
			errs = append(errs, DAYS[hour.DayOfWeek])
		}
	}

	if len(errs) > 0 {
		return status.Errorf(codes.InvalidArgument, fmt.Sprintf("Store with is24hours=true cannot have isOpen=false: %v", strings.Join(errs, ", ")))
	}

	return nil
}

func (g *GrpcRoute) UpdateProductType(ctx context.Context, req *pb.UpsertProductTypeRequest) (*pb.UpsertProductTypeResponse, error) {

	prodTypeId := req.GetId()
	if prodTypeId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	req.ProductType.Id = prodTypeId

	prodType, err := g.validateUpdateProductTypeRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := g.service.UpsertProductType(ctx, prodType); err != nil {
		return nil, err
	}

	return &pb.UpsertProductTypeResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}, nil
}

func (g *GrpcRoute) validateUpdateProductTypeRequest(ctx context.Context, req *pb.UpsertProductTypeRequest) (*prodEntity.ProductType, error) {
	if req.ProductType.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	if req.ProductType.ProductCategoryId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "product category id is required")
	}

	prodType := prodEntity.ProductType{}
	if err := prodType.FromProto(req.ProductType); err != nil {
		return nil, err
	}

	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Error when getting claims from jwt token")
	}
	prodType.CreatedBy = claims.UserID

	return &prodType, nil
}
