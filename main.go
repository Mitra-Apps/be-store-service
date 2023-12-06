package main

import (
	"log"
	"net"

	"github.com/Mitra-Apps/be-store-service/config/postgre"
	pb "github.com/Mitra-Apps/be-store-service/domain/proto/store"
	storePostgreRepo "github.com/Mitra-Apps/be-store-service/domain/store/repository/postgre"
	grpcRoute "github.com/Mitra-Apps/be-store-service/handler/grpc"
	"github.com/Mitra-Apps/be-store-service/lib"
	"github.com/Mitra-Apps/be-store-service/service"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedStoreServiceServer
}

func main() {
	envInit()

	lis, err := net.Listen("tcp", lib.GetEnv("APP_PORT"))
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

func envInit() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
}
