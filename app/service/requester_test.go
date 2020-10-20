package service_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/isac-caja/circuit-breaker-cases/app/service"
	"github.com/isac-caja/circuit-breaker-cases/app/service/mocks"
	"github.com/stretchr/testify/assert"
)

func createMocks() (context.Context, *service.RequesterConfig) {
	ctx := context.Background()
	config := new(service.RequesterConfig)
	return ctx, config
}

func TestSyncRequester(t *testing.T) {

	t.Run("Requester does not create Any client and returns 200 for all", func(t *testing.T) {
		ctx, config := createMocks()
		requester := service.NewServiceRequester(config, nil)
		status, err := requester.Request(ctx)
		assert.Equal(t, http.StatusOK, status)
		assert.Nil(t, err)
	})

	t.Run("Requester must create 1 client and request it.", func(t *testing.T) {
		var clients service.Clients
		ctx, config := createMocks()

		config.ClientConfigList = service.ClientConfigList{
			{Endpoint: "alpha.com.br", Timeout: 100 * time.Millisecond},
		}

		client := new(mocks.Client)
		client.
			On("Get", mock.Anything).
			Return(http.StatusOK, nil)

		clients = append(clients, client)

		requester := service.NewServiceRequester(config, clients)
		status, err := requester.Request(ctx)
		assert.Equal(t, http.StatusOK, status)
		assert.Nil(t, err)
	})

	type RequestCase struct {
		ClientReturnStatuses []int
		ResultStatus         int
	}

	testRequestCase := func(t *testing.T, sCase *RequestCase) {
		var clients service.Clients
		ctx, config := createMocks()

		for _, status := range sCase.ClientReturnStatuses {
			config.ClientConfigList = append(
				config.ClientConfigList,
				service.ClientConfig{Endpoint: "alpha.com.br", Timeout: 100 * time.Millisecond},
			)
			client := new(mocks.Client)
			client.
				On("Get", mock.Anything).
				Return(status, nil)
			clients = append(clients, client)
		}

		requester := service.NewServiceRequester(config, clients)
		status, err := requester.Request(ctx)
		assert.Equal(t, sCase.ResultStatus, status)
		assert.Nil(t, err)
	}

	t.Run("Requester must make 2 successful requests", func(t *testing.T) {
		testRequestCase(t, &RequestCase{
			[]int{http.StatusOK, http.StatusOK},
			http.StatusOK,
		})
	})

	t.Run("Requester must make 2 requests with 1 fails and return 500", func(t *testing.T) {
		testRequestCase(t, &RequestCase{
			[]int{http.StatusOK, http.StatusInternalServerError},
			http.StatusInternalServerError,
		})
	})

	t.Run("Requester must make 2 requests with 1 timeout and return 500", func(t *testing.T) {
		testRequestCase(t, &RequestCase{
			[]int{http.StatusOK, http.StatusRequestTimeout},
			http.StatusInternalServerError,
		})
	})
}

// func TestAsyncRequester(t *testing.T) {

// 	t.Run("Should requests concorrently", func(t *testing.T) {
// 		var clients service.Clients
// 		ctx, config := createMocks()
// 		config.RequesterMode = service.Async

// 		processTime := 1 * time.Second
// 		clientA := new(mocks.Client)
// 		clientA.
// 			On("Get", mock.Anything).
// 			Once().
// 			After(processTime).
// 			Return(http.StatusOK, nil)

// 		clientB := new(mocks.Client)
// 		clientB.
// 			On("Get", mock.Anything).
// 			Once().
// 			After(processTime).
// 			Return(http.StatusOK, nil)

// 		clients = append(clients, clientA)
// 		clients = append(clients, clientB)

// 		requester := service.NewServiceRequester(config, clients)
// 		started := time.Now()
// 		status, err := requester.Request(ctx)
// 		duration := int(time.Since(started).Seconds())

// 		assert.Nil(t, err)
// 		assert.Equal(t, http.StatusOK, status)
// 		assert.Equal(t, int(processTime.Seconds()), duration)

// 		clientA.AssertExpectations(t)
// 		clientB.AssertExpectations(t)
// 	})

// 	t.Run("Should fails 1 request", func(t *testing.T) {
// 		var clients service.Clients
// 		ctx, config := createMocks()
// 		config.RequesterMode = service.Async

// 		processTime := 1 * time.Second
// 		clientA := new(mocks.Client)
// 		clientA.
// 			On("Get", mock.Anything).
// 			Once().
// 			After(processTime).
// 			Return(0, errors.New("Fail"))

// 		clientB := new(mocks.Client)
// 		clientB.
// 			On("Get", mock.Anything).
// 			Once().
// 			After(processTime).
// 			Return(http.StatusOK, nil)

// 		clients = append(clients, clientA)
// 		clients = append(clients, clientB)

// 		requester := service.NewServiceRequester(config, clients)
// 		status, err := requester.Request(ctx)

// 		assert.Equal(t, errors.New("Fail"), err)
// 		assert.Equal(t, 0, status)

// 		clientA.AssertExpectations(t)
// 		clientB.AssertExpectations(t)
// 	})
// }
