package server

import (
	"context"
	"fransfer/internal/auth"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WithStreamServerAuthorizationInterceptor(jwt auth.JWT) grpc.ServerOption {
	return grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(authorize(jwt)))
}

func authorize(jwt auth.JWT) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		token := metautils.ExtractIncoming(ctx).Get(auth.HeaderApiKey)

		if isValid, err := jwt.Validate(token); err != nil || !isValid {
			return nil, status.Error(codes.Unauthenticated, "jwt not valid")
		}
		return context.Background(), nil
	}
}
