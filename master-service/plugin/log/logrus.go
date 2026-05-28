package logrusWrapper

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
)

// NewHandlerWrapper accepts a logrus logging and returns a Handler Wrapper
func NewHandlerWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		start := time.Now()

		l := logrus.New()
		bd, _ := json.Marshal(req.Body())
		rs, _ := json.Marshal(rsp)

		l.Out = os.Stdout

		l.WithFields(logrus.Fields{
			"service":    req.Service(),
			"endpoint":   req.Endpoint(),
			"method":     req.Method(),
			"latency_ns": time.Since(start).Nanoseconds(),
			"body":       string(bd),
			"is_stream":  req.Stream(),
			"response":   string(rs),
		}).Info("[wrapper]")
		err := fn(ctx, req, rsp)
		return err
	}
}

// NewCallWrapper accepts a logrus logging and returns a Call Wrapper
func NewCallWrapper() client.CallWrapper {
	return func(cf client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
			start := time.Now()

			l := logrus.New()
			l.WithFields(logrus.Fields{
				"service":    req.Service(),
				"endpoint":   req.Endpoint(),
				"method":     req.Method(),
				"latency_ns": time.Since(start).Nanoseconds(),
				"body":       req.Body(),
				"is_stream":  req.Stream(),
				"response":   rsp,
			}).Info("request details")

			return nil
		}
	}
}
