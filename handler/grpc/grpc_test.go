package grpc

import (
	"context"
	"testing"

	pb "github.com/Mitra-Apps/be-store-service/domain/proto/store"
	serviceMock "github.com/Mitra-Apps/be-store-service/service/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	userID  = "8b15140c-f6d0-4f2f-8302-57383a51adaf"
	storeID = "7d56be32-70a2-4f49-b66b-63e6f8e719d5"
)

func TestGrpcRoute_OpenCloseStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	svcMock := serviceMock.NewMockService(ctrl)
	ctx := context.Background()
	svcMock.EXPECT().OpenCloseStore(gomock.Any(), gomock.Any(), gomock.Any(), storeID, false).
		Return(nil)
	svcMock.EXPECT().OpenCloseStore(gomock.Any(), gomock.Any(), gomock.Any(), "", false).
		Return(status.Errorf(codes.InvalidArgument, codes.InvalidArgument.String()))

	type args struct {
		ctx context.Context
		req *pb.OpenCloseStoreRequest
	}
	tests := []struct {
		name    string
		s       *GrpcRoute
		args    args
		want    *pb.OpenCloseStoreResponse
		wantErr bool
	}{
		{
			name: "OpenCloseStore_NoError_Success",
			s: &GrpcRoute{
				service: svcMock,
			},
			args: args{
				ctx: ctx,
				req: &pb.OpenCloseStoreRequest{
					StoreId:  storeID,
					IsActive: false,
				},
			},
			want: &pb.OpenCloseStoreResponse{
				Code:    int32(codes.OK),
				Message: codes.OK.String(),
			},
			wantErr: false,
		},
		{
			name: "OpenCloseStore_NoStoreID_InvalidArgumentError",
			s: &GrpcRoute{
				service: svcMock,
			},
			args: args{
				ctx: ctx,
				req: &pb.OpenCloseStoreRequest{
					StoreId:  "",
					IsActive: false,
				},
			},
			want: &pb.OpenCloseStoreResponse{
				Code:    int32(codes.InvalidArgument),
				Message: codes.InvalidArgument.String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.OpenCloseStore(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}
