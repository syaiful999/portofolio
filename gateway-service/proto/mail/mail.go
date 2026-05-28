package masterdata

import (
	"context"
	"moyo-gateway-service/config"
	mail "moyo-gateway-service/proto/mail/mail"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func Register(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption, conf config.Config) error {
	var err error
	mailEndpoint := conf.Services.MailURL

	if err = mail.RegisterMailServiceHandlerFromEndpoint(ctx, mux, mailEndpoint, opts); err != nil {
		return err
	}
	return nil

}
