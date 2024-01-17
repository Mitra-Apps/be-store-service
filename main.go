package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	pb "github.com/Mitra-Apps/be-store-service/domain/proto/store"
	grpcRoute "github.com/Mitra-Apps/be-store-service/handler/grpc"
	"github.com/Mitra-Apps/be-store-service/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmgrpc"

	configPostgres "github.com/Mitra-Apps/be-store-service/config/postgres"
	repositoryPostgres "github.com/Mitra-Apps/be-store-service/domain/store/repository/postgres"
	"github.com/Mitra-Apps/be-store-service/domain/store/repository/storage"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	godotenv.Load()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("GRPC_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db := configPostgres.Connection()
	repoPostgres := repositoryPostgres.NewPostgres(db)
	repoStorage := storage.New()
	svc := service.New(repoPostgres, repoStorage)
	grpcServer := GrpcNewServer(ctx, []grpc.ServerOption{})
	route := grpcRoute.New(svc)
	pb.RegisterStoreServiceServer(grpcServer, route)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	go HttpNewServer(ctx, os.Getenv("GRPC_PORT"), os.Getenv("HTTP_PORT"))

	grpcServer.Serve(lis)
}

func GrpcNewServer(ctx context.Context, opts []grpc.ServerOption) *grpc.Server {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	logrusOpts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
	}
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	opts = append(opts, grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrusEntry, logrusOpts...),
			grpc_recovery.StreamServerInterceptor(),
			apmgrpc.NewStreamServerInterceptor(apmgrpc.WithRecovery()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrusEntry, logrusOpts...),
			grpc_recovery.UnaryServerInterceptor(),
			apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery()),
		)),
	)

	myServer := grpc.NewServer(opts...)

	reflection.Register(myServer)
	return myServer
}

func HttpNewServer(ctx context.Context, grpcPort, httpPort string) error {
	mux := runtime.NewServeMux()
	mux.HandlePath("GET", "/docs/v1/stores/openapi.yaml", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		http.ServeFile(w, r, "docs/openapi.yaml")
	})

	mux.HandlePath("GET", "/docs/v1/stores", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		http.ServeFile(w, r, "docs/index.html")
	})

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterStoreServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%s", grpcPort), opts); err != nil {
		return err
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", httpPort),
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			logrus.Panicf("failed to shutdown server: %v", err)
		}
	}()

	return srv.ListenAndServe()
}
