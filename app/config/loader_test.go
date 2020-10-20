package config_test

import (
	"testing"
	"time"

	"github.com/isac-caja/circuit-breaker-cases/app/config"
	"github.com/isac-caja/circuit-breaker-cases/app/service"
	"github.com/stretchr/testify/assert"
)

const MockConfigPath = "./mocks/config.mock.yaml"

func TestLoader(t *testing.T) {

	t.Run("Load the config file", func(t *testing.T) {
		expected := &service.Config{
			WaitTime:     1 * time.Millisecond,
			ReturnStatus: 2,
			RequesterConfig: service.RequesterConfig{
				RequesterMode: service.Async,
				ClientConfigList: service.ClientConfigList{
					{
						Endpoint: "http://alpha.com/beta",
						Timeout:  3,
						CircuitBreakerConfig: service.CircuitBreakerConfig{
							CircuitBreakerStrategy: service.Fixed,
							AllowFailures:          20,
						},
					},
					{
						Endpoint: "http://charlie.com/delta",
						Timeout:  4,
						CircuitBreakerConfig: service.CircuitBreakerConfig{
							CircuitBreakerStrategy: service.Percentage,
							Factor:                 0.3,
							MaxRequests:            10,
							Interval:               200 * time.Second,
							Timeout:                20 * time.Second,
						},
					},
				},
			},
		}

		loader := config.NewYamlLoader(MockConfigPath)
		config, err := loader.Load()

		assert.Nil(t, err)
		assert.Equal(t, expected, config)
	})
}
