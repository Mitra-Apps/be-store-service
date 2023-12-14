package main

import (
	"log"
	"net"
	"os"

	"github.com/Mitra-Apps/be-store-service/config/postgre"
	pb "github.com/Mitra-Apps/be-store-service/domain/proto/store"
	storePostgreRepo "github.com/Mitra-Apps/be-store-service/domain/store/repository/postgre"
	grpcRoute "github.com/Mitra-Apps/be-store-service/handler/grpc"
	"github.com/Mitra-Apps/be-store-service/service"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db := postgre.Connection()
	usrRepo := storePostgreRepo.NewPostgre(db)
	svc := service.New(usrRepo)
	grpcServer := grpc.NewServer()
	route := grpcRoute.New(svc)
	pb.RegisterStoreServiceServer(grpcServer, route)

	log.Printf("GRPC Server listening at %v ", lis.Addr())
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v \n", err)
	}
}
