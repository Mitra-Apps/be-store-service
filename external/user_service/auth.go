package user_service

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Mitra-Apps/be-store-service/types"
	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
)

type ValidateUserTokenRequest = pb.ValidateUserTokenRequest

type authClient struct {
	client   pb.UserServiceClient
	grpcConn *grpc.ClientConn
}

//go:generate mockgen -source=auth.go -destination=mock/auth.go -package=mock
type Authentication interface {
	Close()
	ValidateUserToken(ctx context.Context, req *emptypb.Empty) (*types.JwtCustomClaim, error)
}

// Authentication client constructor
func NewAuthClient(ctx context.Context) *authClient {

	host := os.Getenv("GRPC_USER_HOST")
	grpcConn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Cannot connect to utility grpc server ", err)
	}
	client := pb.NewUserServiceClient(grpcConn)

	return &authClient{
		client:   client,
		grpcConn: grpcConn,
	}
}

// New method to close the grpc connection
func (s *authClient) Close() {
	log.Println("Closing connection ...")
	if err := s.grpcConn.Close(); err != nil {
		log.Printf("Error closing the connection: %v", err)
	}
}

func (r *authClient) ValidateUserToken(ctx context.Context, req *ValidateUserTokenRequest) (response *types.JwtCustomClaim, err error) {
	resp, err := r.client.ValidateUserToken(ctx, req)
	if err != nil {
		return
	}

	marsh, _ := resp.Data.MarshalJSON()
	if err := json.Unmarshal(marsh, &response); err != nil {
		return nil, err
	}

	return
}
