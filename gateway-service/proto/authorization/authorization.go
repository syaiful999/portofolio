package authorization

import (
	"context"
	"moyo-gateway-service/config"
	auth "moyo-gateway-service/proto/authorization/auth"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func Register(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption, conf config.Config) error {
	var err error
	authEndpoint := conf.Services.AuthURL

	if err = auth.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, authEndpoint, opts); err != nil {
		return err
	}

	return nil

}
