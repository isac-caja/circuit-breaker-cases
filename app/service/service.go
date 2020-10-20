//go:generate mockery --name=ServiceRequester --case=snake

package service

import (
	"context"
	"log"
	"time"
)

type ServiceRequester interface {
	Request(context.Context) (int, error)
}

type (
	Config struct {
		WaitTime     time.Duration
		ReturnStatus int
		RequesterConfig
	}

	Service struct {
		*Config
		ServiceRequester
	}
)

func New(config *Config, requester ServiceRequester) *Service {
	svc := new(Service)
	svc.Config = config
	svc.ServiceRequester = requester
	return svc
}

func (svc *Service) Process(ctx context.Context) (int, error) {
	time.Sleep(svc.WaitTime)
	if svc.ReturnStatus > 0 {
		return svc.ReturnStatus, nil
	}
	return svc.ServiceRequester.Request(ctx)
}

func addProcessLog(ctx context.Context) {
	if v := ctx.Value("logger"); v != nil {
		logger := v.(*log.Logger)
		logger.Println("Processing")
	}
}
