// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: proto/store/store.proto

package storeconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	store "github.com/Mitra-Apps/be-store-service/domain/proto/store"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// StoreServiceName is the fully-qualified name of the StoreService service.
	StoreServiceName = "StoreService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// StoreServiceCreateStoreProcedure is the fully-qualified name of the StoreService's CreateStore
	// RPC.
	StoreServiceCreateStoreProcedure = "/StoreService/CreateStore"
	// StoreServiceGetStoreProcedure is the fully-qualified name of the StoreService's GetStore RPC.
	StoreServiceGetStoreProcedure = "/StoreService/GetStore"
	// StoreServiceUpdateStoreProcedure is the fully-qualified name of the StoreService's UpdateStore
	// RPC.
	StoreServiceUpdateStoreProcedure = "/StoreService/UpdateStore"
	// StoreServiceDeleteStoreProcedure is the fully-qualified name of the StoreService's DeleteStore
	// RPC.
	StoreServiceDeleteStoreProcedure = "/StoreService/DeleteStore"
	// StoreServiceListStoresProcedure is the fully-qualified name of the StoreService's ListStores RPC.
	StoreServiceListStoresProcedure = "/StoreService/ListStores"
	// StoreServiceOpenCloseStoreProcedure is the fully-qualified name of the StoreService's
	// OpenCloseStore RPC.
	StoreServiceOpenCloseStoreProcedure = "/StoreService/OpenCloseStore"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	storeServiceServiceDescriptor              = store.File_proto_store_store_proto.Services().ByName("StoreService")
	storeServiceCreateStoreMethodDescriptor    = storeServiceServiceDescriptor.Methods().ByName("CreateStore")
	storeServiceGetStoreMethodDescriptor       = storeServiceServiceDescriptor.Methods().ByName("GetStore")
	storeServiceUpdateStoreMethodDescriptor    = storeServiceServiceDescriptor.Methods().ByName("UpdateStore")
	storeServiceDeleteStoreMethodDescriptor    = storeServiceServiceDescriptor.Methods().ByName("DeleteStore")
	storeServiceListStoresMethodDescriptor     = storeServiceServiceDescriptor.Methods().ByName("ListStores")
	storeServiceOpenCloseStoreMethodDescriptor = storeServiceServiceDescriptor.Methods().ByName("OpenCloseStore")
)

// StoreServiceClient is a client for the StoreService service.
type StoreServiceClient interface {
	// Create a new store
	CreateStore(context.Context, *connect.Request[store.CreateStoreRequest]) (*connect.Response[store.Store], error)
	// Get a store by ID
	GetStore(context.Context, *connect.Request[store.GetStoreRequest]) (*connect.Response[store.Store], error)
	// Update an existing store
	UpdateStore(context.Context, *connect.Request[store.UpdateStoreRequest]) (*connect.Response[store.Store], error)
	// Delete a store by ID
	DeleteStore(context.Context, *connect.Request[store.DeleteStoreRequest]) (*connect.Response[emptypb.Empty], error)
	// List all stores
	ListStores(context.Context, *connect.Request[store.ListStoresRequest]) (*connect.Response[store.ListStoresResponse], error)
	// Open close store
	OpenCloseStore(context.Context, *connect.Request[store.OpenCloseStoreRequest]) (*connect.Response[emptypb.Empty], error)
}

// NewStoreServiceClient constructs a client for the StoreService service. By default, it uses the
// Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewStoreServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) StoreServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &storeServiceClient{
		createStore: connect.NewClient[store.CreateStoreRequest, store.Store](
			httpClient,
			baseURL+StoreServiceCreateStoreProcedure,
			connect.WithSchema(storeServiceCreateStoreMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getStore: connect.NewClient[store.GetStoreRequest, store.Store](
			httpClient,
			baseURL+StoreServiceGetStoreProcedure,
			connect.WithSchema(storeServiceGetStoreMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updateStore: connect.NewClient[store.UpdateStoreRequest, store.Store](
			httpClient,
			baseURL+StoreServiceUpdateStoreProcedure,
			connect.WithSchema(storeServiceUpdateStoreMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		deleteStore: connect.NewClient[store.DeleteStoreRequest, emptypb.Empty](
			httpClient,
			baseURL+StoreServiceDeleteStoreProcedure,
			connect.WithSchema(storeServiceDeleteStoreMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listStores: connect.NewClient[store.ListStoresRequest, store.ListStoresResponse](
			httpClient,
			baseURL+StoreServiceListStoresProcedure,
			connect.WithSchema(storeServiceListStoresMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		openCloseStore: connect.NewClient[store.OpenCloseStoreRequest, emptypb.Empty](
			httpClient,
			baseURL+StoreServiceOpenCloseStoreProcedure,
			connect.WithSchema(storeServiceOpenCloseStoreMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// storeServiceClient implements StoreServiceClient.
type storeServiceClient struct {
	createStore    *connect.Client[store.CreateStoreRequest, store.Store]
	getStore       *connect.Client[store.GetStoreRequest, store.Store]
	updateStore    *connect.Client[store.UpdateStoreRequest, store.Store]
	deleteStore    *connect.Client[store.DeleteStoreRequest, emptypb.Empty]
	listStores     *connect.Client[store.ListStoresRequest, store.ListStoresResponse]
	openCloseStore *connect.Client[store.OpenCloseStoreRequest, emptypb.Empty]
}

// CreateStore calls StoreService.CreateStore.
func (c *storeServiceClient) CreateStore(ctx context.Context, req *connect.Request[store.CreateStoreRequest]) (*connect.Response[store.Store], error) {
	return c.createStore.CallUnary(ctx, req)
}

// GetStore calls StoreService.GetStore.
func (c *storeServiceClient) GetStore(ctx context.Context, req *connect.Request[store.GetStoreRequest]) (*connect.Response[store.Store], error) {
	return c.getStore.CallUnary(ctx, req)
}

// UpdateStore calls StoreService.UpdateStore.
func (c *storeServiceClient) UpdateStore(ctx context.Context, req *connect.Request[store.UpdateStoreRequest]) (*connect.Response[store.Store], error) {
	return c.updateStore.CallUnary(ctx, req)
}

// DeleteStore calls StoreService.DeleteStore.
func (c *storeServiceClient) DeleteStore(ctx context.Context, req *connect.Request[store.DeleteStoreRequest]) (*connect.Response[emptypb.Empty], error) {
	return c.deleteStore.CallUnary(ctx, req)
}

// ListStores calls StoreService.ListStores.
func (c *storeServiceClient) ListStores(ctx context.Context, req *connect.Request[store.ListStoresRequest]) (*connect.Response[store.ListStoresResponse], error) {
	return c.listStores.CallUnary(ctx, req)
}

// OpenCloseStore calls StoreService.OpenCloseStore.
func (c *storeServiceClient) OpenCloseStore(ctx context.Context, req *connect.Request[store.OpenCloseStoreRequest]) (*connect.Response[emptypb.Empty], error) {
	return c.openCloseStore.CallUnary(ctx, req)
}

// StoreServiceHandler is an implementation of the StoreService service.
type StoreServiceHandler interface {
	// Create a new store
	CreateStore(context.Context, *connect.Request[store.CreateStoreRequest]) (*connect.Response[store.Store], error)
	// Get a store by ID
	GetStore(context.Context, *connect.Request[store.GetStoreRequest]) (*connect.Response[store.Store], error)
	// Update an existing store
	UpdateStore(context.Context, *connect.Request[store.UpdateStoreRequest]) (*connect.Response[store.Store], error)
	// Delete a store by ID
	DeleteStore(context.Context, *connect.Request[store.DeleteStoreRequest]) (*connect.Response[emptypb.Empty], error)
	// List all stores
	ListStores(context.Context, *connect.Request[store.ListStoresRequest]) (*connect.Response[store.ListStoresResponse], error)
	// Open close store
	OpenCloseStore(context.Context, *connect.Request[store.OpenCloseStoreRequest]) (*connect.Response[emptypb.Empty], error)
}

// NewStoreServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewStoreServiceHandler(svc StoreServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	storeServiceCreateStoreHandler := connect.NewUnaryHandler(
		StoreServiceCreateStoreProcedure,
		svc.CreateStore,
		connect.WithSchema(storeServiceCreateStoreMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	storeServiceGetStoreHandler := connect.NewUnaryHandler(
		StoreServiceGetStoreProcedure,
		svc.GetStore,
		connect.WithSchema(storeServiceGetStoreMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	storeServiceUpdateStoreHandler := connect.NewUnaryHandler(
		StoreServiceUpdateStoreProcedure,
		svc.UpdateStore,
		connect.WithSchema(storeServiceUpdateStoreMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	storeServiceDeleteStoreHandler := connect.NewUnaryHandler(
		StoreServiceDeleteStoreProcedure,
		svc.DeleteStore,
		connect.WithSchema(storeServiceDeleteStoreMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	storeServiceListStoresHandler := connect.NewUnaryHandler(
		StoreServiceListStoresProcedure,
		svc.ListStores,
		connect.WithSchema(storeServiceListStoresMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	storeServiceOpenCloseStoreHandler := connect.NewUnaryHandler(
		StoreServiceOpenCloseStoreProcedure,
		svc.OpenCloseStore,
		connect.WithSchema(storeServiceOpenCloseStoreMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/StoreService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case StoreServiceCreateStoreProcedure:
			storeServiceCreateStoreHandler.ServeHTTP(w, r)
		case StoreServiceGetStoreProcedure:
			storeServiceGetStoreHandler.ServeHTTP(w, r)
		case StoreServiceUpdateStoreProcedure:
			storeServiceUpdateStoreHandler.ServeHTTP(w, r)
		case StoreServiceDeleteStoreProcedure:
			storeServiceDeleteStoreHandler.ServeHTTP(w, r)
		case StoreServiceListStoresProcedure:
			storeServiceListStoresHandler.ServeHTTP(w, r)
		case StoreServiceOpenCloseStoreProcedure:
			storeServiceOpenCloseStoreHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedStoreServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedStoreServiceHandler struct{}

func (UnimplementedStoreServiceHandler) CreateStore(context.Context, *connect.Request[store.CreateStoreRequest]) (*connect.Response[store.Store], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("StoreService.CreateStore is not implemented"))
}

func (UnimplementedStoreServiceHandler) GetStore(context.Context, *connect.Request[store.GetStoreRequest]) (*connect.Response[store.Store], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("StoreService.GetStore is not implemented"))
}

func (UnimplementedStoreServiceHandler) UpdateStore(context.Context, *connect.Request[store.UpdateStoreRequest]) (*connect.Response[store.Store], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("StoreService.UpdateStore is not implemented"))
}

func (UnimplementedStoreServiceHandler) DeleteStore(context.Context, *connect.Request[store.DeleteStoreRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("StoreService.DeleteStore is not implemented"))
}

func (UnimplementedStoreServiceHandler) ListStores(context.Context, *connect.Request[store.ListStoresRequest]) (*connect.Response[store.ListStoresResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("StoreService.ListStores is not implemented"))
}

func (UnimplementedStoreServiceHandler) OpenCloseStore(context.Context, *connect.Request[store.OpenCloseStoreRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("StoreService.OpenCloseStore is not implemented"))
}
