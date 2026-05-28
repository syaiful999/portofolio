package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	ch "moyo-gateway-service/cache"
	cf "moyo-gateway-service/config"
	"moyo-gateway-service/middleware"
	"moyo-gateway-service/pkg/entities"
	pkglogs "moyo-gateway-service/pkg/logs"
	master_data "moyo-gateway-service/proto/master-data"
	"moyo-gateway-service/utils"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	corsMux "github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	conf  = cf.Config{}
	d     struct{}
	cache ch.Cache
)

func run() error {
	var cancel context.CancelFunc
	utils.PushLogf("", "START", "------------------------------------------------------------------------------------")
	if err := godotenv.Load(); err != nil {
		return err
	}

	conf.Init()

	ctx := context.Background()
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	cache, err := ch.Connection(&conf)
	if err != nil {
		return err
	}

	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			header := request.Header.Get("Authorization")
			clientId := request.Header.Get("x-client-id")
			clientSecret := request.Header.Get("x-client-secret")
			realm := request.Header.Get("x-realms")
			userToken := request.Header.Get("x-user-token")
			checkLog := request.Header.Get("x-check-log")

			var headers entities.Headers
			generateHeaders(&headers, request, cache)

			md := metadata.Pairs(
				"Authorization", header,
				"x-client-id", clientId,
				"x-client-secret", clientSecret,
				"x-realms", realm,
				"x-user-id", headers.X_User_ID,
				"x-user-email", headers.X_User_Email,
				"x-user-fullname", headers.X_User_Fullname,
				"x-company-id", headers.X_Company_ID,
				"x-company-name", headers.X_Company_Name,
				"x-user-token", userToken,
				"x-check-log", checkLog,
			)

			return md
		}),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := master_data.Register(ctx, mux, opts, conf); err != nil {
		return err
	}

	server := fiber.New(
		fiber.Config{
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"Status":    fiber.StatusInternalServerError,
					"IsSuccess": false,
					"Message":   err.Error(),
					"Data":      d,
				})
			},
			AppName: conf.Service.Name,
		},
	)
	Handler(server)

	go func() {
		portFiber := "8081"
		if conf.Service.PortFiber != 0 {
			portFiber = strconv.Itoa(conf.Service.PortFiber)
		}
		if err := server.Listen(":" + portFiber); err != nil {
			glog.Fatalf("Fiber failed to start: %v", err)
		}
	}()

	portHttp := "8080"
	if conf.Service.Port != 0 {
		portHttp = strconv.Itoa(conf.Service.Port)
	}
	gatewayMux := http.NewServeMux()
	gatewayMux.Handle("/api/v1/", mux)

	corsMiddleware := corsMux.New(corsMux.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-User-Token", "Accept", "Accept-Encoding", "ResponseType", "Referrer-Policy"},
		AllowCredentials: false,
	})

	handler := corsMiddleware.Handler(gatewayMux)

	return http.ListenAndServe(":"+portHttp, handler)
}

func Handler(route *fiber.App) {
	route.Use(middleware.CustomLogger())

	route.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Accept,Content-Type,Content-Length,Accept-Encoding,Authorization,ResponseType,X-User-Token,Referrer-Policy",
		AllowCredentials: false,
	}))

	route.Get("/health-check", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"IsSuccess": true,
			"Message":   fmt.Sprintf("SVC RUNNING %s - Version %s", conf.Service.Name, conf.Service.Version),
			"Data":      d,
			"Status":    "0",
		})
	})

	route.Get("/logs/:date", pkglogs.Handler)

	route.Get("/access", middleware.Protected(), func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"IsSuccess": true,
			"Message":   "Grant Access",
			"Data":      d,
			"Status":    "0",
		})
	})
}

func main() {
	flag.Parse()

	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

func generateHeaders(h *entities.Headers, request *http.Request, cache ch.Cache) {
	authorization := request.Header.Get("Authorization")
	sub, ok := utils.GetSub(authorization)

	if !ok {
		return
	}

	data, _, err := cache.Get(sub)
	if err != nil {
		return
	}

	var rbacEnt entities.RBAC
	if err := json.Unmarshal(data.([]byte), &rbacEnt); err != nil {
		return
	}

	h.Generate(
		rbacEnt.User.User.ID,
		rbacEnt.User.Account.Email,
		fmt.Sprintf("%s %s", rbacEnt.User.User.FirstName, rbacEnt.User.User.LastName),
		rbacEnt.User.Company.ID,
		rbacEnt.User.Company.Name,
	)
}
