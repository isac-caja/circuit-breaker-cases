package service_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/isac-caja/circuit-breaker-cases/app/service"
	"github.com/isac-caja/circuit-breaker-cases/app/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {

	createMocks := func() (context.Context, *service.Config, *mocks.ServiceRequester) {
		ctx := context.Background()
		config := new(service.Config)
		req := new(mocks.ServiceRequester)
		return ctx, config, req
	}

	t.Run("Test a simple process call: success", func(t *testing.T) {
		ctx, config, requester := createMocks()
		svc := service.New(config, requester)

		requester.
			On("Request", mock.Anything).
			Once().
			Return(http.StatusOK, nil)

		status, err := svc.Process(ctx)
		assert.Equal(t, http.StatusOK, status)
		assert.Nil(t, err)
	})

	t.Run("Must return an error of requester", func(t *testing.T) {
		ctx, config, requester := createMocks()
		svc := service.New(config, requester)

		requester.
			On("Request", mock.Anything).
			Once().
			Return(0, errors.New("Testing"))

		status, err := svc.Process(ctx)
		assert.Equal(t, 0, status)
		assert.Equal(t, "Testing", err.Error())
	})

	t.Run("Must wait for 1s", func(t *testing.T) {
		ctx, config, requester := createMocks()
		svc := service.New(config, requester)

		config.WaitTime = time.Duration(1 * time.Second)

		requester.
			On("Request", mock.Anything).
			Once().
			Return(http.StatusOK, nil)

		start := time.Now()
		status, err := svc.Process(ctx)
		assert.LessOrEqual(t, float64(1.0), time.Since(start).Seconds())
		assert.Equal(t, http.StatusOK, status)
		assert.Nil(t, err)
	})

	t.Run("Must return the status 400", func(t *testing.T) {
		ctx, config, requester := createMocks()
		svc := service.New(config, requester)

		config.WaitTime = time.Duration(1 * time.Second)
		config.ReturnStatus = http.StatusBadRequest

		start := time.Now()
		status, err := svc.Process(ctx)
		assert.LessOrEqual(t, float64(1.0), time.Since(start).Seconds())
		assert.Equal(t, http.StatusBadRequest, status)
		assert.Nil(t, err)
	})
}
