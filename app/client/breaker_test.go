package client_test

import (
	"context"
	"errors"
	"testing"

	"github.com/isac-caja/circuit-breaker-cases/app/client"
	"github.com/isac-caja/circuit-breaker-cases/app/service"
	"github.com/isac-caja/circuit-breaker-cases/app/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFixedCircuitBreaker(t *testing.T) {

	t.Run("Should open the circuit with 10 requests", func(t *testing.T) {
		ctx := context.Background()
		config := service.CircuitBreakerConfig{
			AllowFailures:          10,
			CircuitBreakerStrategy: service.Fixed,
		}
		mockclient := new(mocks.Client)
		mockclient.
			On("Get", mock.Anything).
			Times(11).
			Return(0, errors.New("Fail"))

		clBreaker := client.WithBreaker(mockclient, client.NewBreaker(config))

		var i int
		for i <= 20 {
			status, err := clBreaker.Get(ctx)
			assert.Equal(t, 0, status)
			assert.Error(t, err, "An error was expected at %dth interaction", i)
			i++
		}
		mockclient.AssertExpectations(t)
	})

	t.Run("Should open the circuit at 101th request", func(t *testing.T) {
		ctx := context.Background()
		config := service.CircuitBreakerConfig{
			Factor:                 0.3,
			CircuitBreakerStrategy: service.Percentage,
		}
		mockclient := new(mocks.Client)
		mockclient.
			On("Get", mock.Anything).
			Times(70).
			Return(200, nil)

		mockclient.
			On("Get", mock.Anything).
			Times(31).
			Return(0, errors.New("Fail"))

		clBreaker := client.WithBreaker(mockclient, client.NewBreaker(config))

		var i int
		for i <= 200 {
			i++
			clBreaker.Get(ctx)
		}
		mockclient.AssertExpectations(t)
	})

}
