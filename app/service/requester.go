//go:generate mockery --name=Client --case=snake

package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type (
	Clients []Client
	Client  interface {
		Get(context.Context) (int, error)
	}
)

type (
	RequesterConfig struct {
		ClientConfigList
		RequesterMode
	}

	ClientConfigList []ClientConfig
	ClientConfig     struct {
		Endpoint string
		Timeout  time.Duration
		CircuitBreakerConfig
	}

	CircuitBreakerConfig struct {
		CircuitBreakerStrategy
		MaxRequests   uint32
		Interval      time.Duration
		Timeout       time.Duration
		Factor        float32
		AllowFailures uint32
	}
)

type (
	RequesterMode          int
	CircuitBreakerStrategy int
)

const (
	Sync RequesterMode = iota
	Async

	Fixed CircuitBreakerStrategy = 1 << iota
	Percentage
)

func NewServiceRequester(config *RequesterConfig, clients Clients) ServiceRequester {
	var req ServiceRequester

	if len(clients) < 1 {
		return new(DummyRequester)
	}

	switch config.RequesterMode {
	case Sync:
		req = &SyncRequester{
			clients: clients,
		}
	// case Async:
	// 	req = &AsyncRequester{
	// 		clients: clients,
	// 	}
	default:
		panic(fmt.Errorf("Requester Mode not implemented: %v", config.RequesterMode))
	}
	return req
}

type DummyRequester struct{}

func (dr *DummyRequester) Request(ctx context.Context) (int, error) {
	return http.StatusOK, nil
}

type SyncRequester struct {
	clients Clients
}

func (sr *SyncRequester) Request(ctx context.Context) (int, error) {
	for _, c := range sr.clients {
		status, err := c.Get(ctx)
		if err != nil {
			return 0, err
		}
		if !IsValid(status) {
			return http.StatusInternalServerError, nil
		}
	}
	return http.StatusOK, nil
}

// type AsyncRequester struct {
// 	clients Clients
// }

// func (ar *AsyncRequester) Request(ctx context.Context) (int, error) {
// 	var wg sync.WaitGroup
// 	type Result struct {
// 		status int
// 		err    error
// 	}

// 	ch := make(chan Result)
// 	doRequest := func(ctx context.Context, wg *sync.WaitGroup, c Client) {
// 		status, err := c.Get(ctx)
// 		ch <- Result{status, err}
// 		wg.Done()
// 	}

// 	wg.Wait()

// 	for res := range ch {
// 		if !IsValid(res.status) {
// 			return http.StatusInternalServerError, nil
// 		}
// 	}

// 	return http.StatusOK, nil
// }

func addSyncRequestLog(ctx context.Context, client Client) {
	if v := ctx.Value("logger"); v != nil {
		logger := v.(*log.Logger)
		logger.Println("Requesting:", client)
	}
}

// func addAsynRequestLog(ctx context.Context, client Client) {
// 	if v := ctx.Value("logger"); v != nil {
// 		logger := v.(*log.Logger)
// 		logger.Println("Requesting Async:", client)
// 	}
// }

func IsValid(status int) bool {
	if status < 300 {
		return true
	}
	return false
}
