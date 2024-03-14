package middleware

import (
	"context"
	"fmt"
	"os"
	"strings"

	userService "github.com/Mitra-Apps/be-user-service/service"
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

func verifyToken(ctx context.Context, tokenString string) (string, []string, error) {
	authClient := userService.NewAuthClient(os.Getenv("JWT_SECRET"))
	claims, err := authClient.ValidateToken(ctx, tokenString)
	if err != nil {
		return "", nil, err
	}

	return claims.Subject, claims.Roles, nil
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
