package middleware

import (
	"context"
	"fmt"
	"log"
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

type JwtClaims struct {
	UserID    uuid.UUID
	RoleNames []string
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
	userId, roleNames, err := verifyToken(token)
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

func verifyToken(tokenString string) (string, []string, error) {
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
		return "", nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil, fmt.Errorf("token claims are not of type jwt.MapClaims")
	}

	userId, userIdOk := claims["sub"].(string)
	rolesRaw, rolesOK := claims["roles"]
	// expirationTime, expOk := claims["exp"].(float64)

	if !userIdOk {
		return "", nil, fmt.Errorf("invalid token")
	}

	var roleNames []string
	if rolesOK {
		roles, ok := rolesRaw.([]interface{})
		if !ok {
			log.Fatal("Error converting roles to []interface{}")
		}

		// Convert each role to a string.
		for _, role := range roles {
			roleString, ok := role.(string)
			if !ok {
				log.Fatal("Error converting role to string")
			}
			roleNames = append(roleNames, roleString)
		}
	}

	return userId, roleNames, nil
}

// GetClaimsFromContext returns the userId from the context
func GetClaimsFromContext(ctx context.Context) (*JwtClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no headers provided")
	}
	userIDVal := md["x-user-id"]
	if len(userIDVal) == 0 {
		return nil, fmt.Errorf("no user id provided")
	}

	userID, err := uuid.Parse(userIDVal[0])
	if err != nil {
		return nil, err
	}
	var jwtClaims JwtClaims
	jwtClaims.UserID = userID

	jwtClaims.RoleNames = md["x-role-names"]

	return &jwtClaims, nil
}
