package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mitra-Apps/be-store-service/external/user_service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var excludedMethods = []string{
	// "/StoreService/GetStore",
}

type JwtClaims struct {
	UserID    uuid.UUID
	RoleNames []string
	IsAdmin   bool
}

func Auth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	for _, excludedMethod := range excludedMethods {
		if strings.Contains(info.FullMethod, excludedMethod) {
			return handler(ctx, req)
		}
	}

	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "no headers provided")
	}

	token := getTokenValue(headers)
	if token == "" {
		return nil, status.Errorf(codes.Unauthenticated, "no token provided")
	}

	// extract the jwt token and get the userId
	userId, roleNames, err := verifyToken(ctx, token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	headers.Append("x-user-id", userId)
	headers.Append("x-role-names", roleNames...)

	ctx = metadata.NewIncomingContext(ctx, headers)

	return handler(ctx, req)
}

func getTokenValue(headers metadata.MD) string {
	token := ""
	if values := headers.Get("authorization"); len(values) > 0 {
		token = strings.TrimPrefix(values[0], "Bearer ")
		token = strings.TrimSpace(token)
	}
	return token
}

func verifyToken(ctx context.Context, token string) (string, []string, error) {
	userService := user_service.NewAuthClient(ctx)
	defer userService.Close()

	params := user_service.ValidateUserTokenRequest{
		Token: token,
	}
	response, err := userService.ValidateUserToken(ctx, &params)

	if err != nil {
		return "", nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return response.RegisteredClaims.Subject, response.Roles, nil
}

func GetClaimsFromContext(ctx context.Context) (*JwtClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no metadata provided")
	}
	userIDStr := md["x-user-id"]
	if len(userIDStr) == 0 {
		return nil, fmt.Errorf("no user ID provided")
	}

	userID, err := uuid.Parse(userIDStr[0])
	if err != nil {
		return nil, err
	}

	var jwtClaims JwtClaims
	jwtClaims.UserID = userID
	jwtClaims.RoleNames = md["x-role-names"]

	isAdmin := false
	for _, role := range jwtClaims.RoleNames {
		if role == "admin" {
			isAdmin = true
		}
	}
	jwtClaims.IsAdmin = isAdmin

	return &jwtClaims, nil
}
