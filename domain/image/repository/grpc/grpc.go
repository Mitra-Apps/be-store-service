package grpc

import (
	"context"

	utilityPb "github.com/Mitra-Apps/be-utility-service/domain/proto/utility"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcClient struct {
	pb utilityPb.ImageServiceClient
}

func New(pb utilityPb.ImageServiceClient) *GrpcClient {
	return &GrpcClient{pb: pb}
}

func (g *GrpcClient) UploadImage(ctx context.Context, imageBase64Str, groupName, userID string) (*uuid.UUID, error) {
	res, err := g.pb.UploadImage(ctx, &utilityPb.UploadImageRequest{
		ImageBase64Str: imageBase64Str,
		UserId:         userID,
		GroupName:      groupName,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, status.Errorf(codes.DataLoss, "Failed to upload image")
	}
	imageID := uuid.MustParse(res.Data.GetId())
	return &imageID, nil
}

func (g *GrpcClient) GetImagesByIds(ctx context.Context, ids []string) ([]*utilityPb.Image, error) {
	res, err := g.pb.GetImagesByIds(ctx, &utilityPb.GetImagesByIdsRequest{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	return res.GetData(), nil
}

func (g *GrpcClient) RemoveImage(ctx context.Context, ids []string, groupName, userID string) error {
	_, err := g.pb.DeleteImages(ctx, &utilityPb.DeleteImagesRequest{
		UserId:    userID,
		GroupName: groupName,
		ImageIds:  ids,
	})
	return err
}
