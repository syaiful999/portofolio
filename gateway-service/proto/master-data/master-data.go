package masterdata

import (
	"context"
	"moyo-gateway-service/config"
	user "moyo-gateway-service/proto/master-data/user"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func Register(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption, conf config.Config) error {
	masterDataEndpoint := conf.Services.MasterDataURL

	if err := user.RegisterUserServiceHandlerFromEndpoint(ctx, mux, masterDataEndpoint, opts); err != nil {
		return err
	}

	return nil
}
