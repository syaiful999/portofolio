package server

import (
	"fmt"
	"log"
	"moyo-master-service/config"
	"moyo-master-service/database"
	enum "moyo-master-service/pkg/enum"
	enumProto "moyo-master-service/pkg/enum/proto"
	user "moyo-master-service/pkg/user"
	userProto "moyo-master-service/pkg/user/proto"
	"strings"

	"github.com/asim/go-micro/plugins/client/grpc/v4"
	"github.com/asim/go-micro/plugins/registry/kubernetes/v4"
	"github.com/asim/go-micro/plugins/registry/mdns/v4"
	gr "github.com/asim/go-micro/plugins/server/grpc/v4"
	_ "github.com/asim/go-micro/plugins/wrapper/monitoring/prometheus/v4"
	_ "github.com/asim/go-micro/plugins/wrapper/monitoring/victoriametrics/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"
)

type Factory struct {
	EnumHandler enum.IEnumHandler
	UserHandler user.IUserHandler
}

func RegisterHandler(srv micro.Service, f Factory) error {
	enumProto.RegisterEnumServiceHandler(srv.Server(), f.EnumHandler)
	userProto.RegisterUserServiceHandler(srv.Server(), f.UserHandler)

	if err := srv.Run(); err != nil {
		return err
	}

	return nil
}

func InitFactory(factory *Factory, conf config.Config) error {
	db := database.DB{
		Name:     conf.Hosts.Database.Name,
		User:     conf.Hosts.Database.Username,
		Password: conf.Hosts.Database.Password,
		Address:  conf.Hosts.Database.Address,
		Port:     conf.Hosts.Database.Port,
		Driver:   conf.Hosts.Database.Driver,
	}

	dbConn, err := database.Connection(&db)
	if err != nil {
		return err
	}

	enumRepository := enum.NewEnumRepository(dbConn)
	userRepository := user.NewUserRepository(dbConn)

	enumUsecase := enum.NewUseCaseEnum(enumRepository)
	userUsecase := user.NewUseCaseUser(userRepository, conf)

	factory.EnumHandler = enum.NewEnumHandler(enumUsecase)
	factory.UserHandler = user.NewUserHandler(userUsecase)

	return nil
}

func Init(conf config.Config) {
	srv := micro.NewService(
		micro.Server(gr.NewServer(server.Name(conf.Service.Name))),
		micro.Client(grpc.NewClient()),
		micro.Name(conf.Service.Name),
		micro.Version(conf.Service.Version),
		micro.Address(fmt.Sprintf("%s:%d", conf.Service.Address, conf.Service.Port)),
		microOptionRegistry(conf),
	)
	srv.Init()

	f := Factory{}
	if err := InitFactory(&f, conf); err != nil {
		log.Fatalf("cannot generate factory: %s", err)
	}

	if err := RegisterHandler(srv, f); err != nil {
		log.Fatalf("register handler error: %s", err)
	}
}

func microOptionRegistry(conf config.Config) micro.Option {
	var registry micro.Option

	switch sd := conf.Hosts.Discovery.Driver; {
	case strings.ToLower(sd) == "kubernetes":
		registry = micro.Registry(
			kubernetes.NewRegistry(),
		)
	default:
		registry = micro.Registry(
			mdns.NewRegistry(),
		)
	}

	return registry
}
