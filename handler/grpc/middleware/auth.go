package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mitra-Apps/be-user-service/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var excludedMethods = []string{
	// "/StoreService/GetStore",
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
	userId, err := verifyToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	headers.Append("x-user-id", userId)

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

func verifyToken(tokenString string) (string, error) {
	// token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
	// 	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	// 	}
	// 	return []byte("secret"), nil
	// })
	// if err != nil {
	// 	return "", err
	// }

	// if !token.Valid {
	// 	return "", fmt.Errorf("invalid token")
	// }

	token, err := auth.VerifyToken(tokenString)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("token claims are not of type jwt.MapClaims")
	}

	userId, userIdOk := claims["userId"].(string)
	// roleNames, roleNamesOk := claims["RoleNames"].([]interface{})
	// expirationTime, expOk := claims["exp"].(float64)

	if !userIdOk {
		return "", fmt.Errorf("invalid token")
	}

	// expiration := time.Unix(int64(expirationTime), 0)
	// fmt.Println("UserID:", userId)
	// fmt.Println("RoleNames:", roleNames)
	// fmt.Println("Expiration Time:", expiration)

	return userId, nil
}

// GetUserIDFromContext returns the userId from the context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, fmt.Errorf("no headers provided")
	}
	values := md["x-user-id"]
	if len(values) == 0 {
		return uuid.Nil, fmt.Errorf("no user id provided")
	}

	userID, err := uuid.Parse(values[0])
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
