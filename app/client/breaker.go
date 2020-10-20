package client

import (
	"context"
	"fmt"

	"github.com/sony/gobreaker"

	"github.com/isac-caja/circuit-breaker-cases/app/service"
)

type (
	Breaker *gobreaker.CircuitBreaker

	ClientBreaker struct {
		service.Client
		Breaker *gobreaker.CircuitBreaker
	}
)

func WithBreaker(client service.Client, breaker Breaker) service.Client {
	if breaker == nil {
		return client
	}
	return &ClientBreaker{
		Client:  client,
		Breaker: breaker,
	}
}

func (cb *ClientBreaker) Get(ctx context.Context) (int, error) {
	process := func() (interface{}, error) {
		status, err := cb.Client.Get(ctx)
		if err != nil {
			return nil, err
		} else if !service.IsValid(status) {
			return nil, fmt.Errorf("Invalid status: %d", status)
		}
		return status, nil
	}

	status, err := cb.Breaker.Execute(process)
	if err != nil {
		return 0, err
	}

	switch v := status.(type) {
	case int:
		return v, nil
	default:
		panic(fmt.Errorf("Inconsistent values status=%#v err=%#v", status, err))
	}
}

func NewBreaker(config service.CircuitBreakerConfig) Breaker {
	switch config.CircuitBreakerStrategy {
	case service.Fixed:
		return NewFixedCircuitBreaker(config)
	case service.Percentage:
		return NewPercentCircuitBreaker(config)
	default:
		return nil
	}
}

func NewFixedCircuitBreaker(config service.CircuitBreakerConfig) Breaker {
	st := newSettings("FixedCircuitBreaker", newFixedValueConditionFunc(config), config)
	return gobreaker.NewCircuitBreaker(st)
}

func NewPercentCircuitBreaker(config service.CircuitBreakerConfig) Breaker {
	st := newSettings("PercentCircuitBreaker", newPercentConditionFunc(config), config)
	return gobreaker.NewCircuitBreaker(st)
}

func newSettings(name string, readyToTrip func(gobreaker.Counts) bool, config service.CircuitBreakerConfig) gobreaker.Settings {
	return gobreaker.Settings{
		Name:          name,
		MaxRequests:   config.MaxRequests,
		Interval:      config.Interval,
		Timeout:       config.Timeout,
		ReadyToTrip:   readyToTrip,
		OnStateChange: onStateChange,
	}
}

func newFixedValueConditionFunc(config service.CircuitBreakerConfig) func(gobreaker.Counts) bool {
	return func(counts gobreaker.Counts) bool {
		if counts.ConsecutiveFailures > config.AllowFailures {
			return true
		}
		return false
	}
}

func newPercentConditionFunc(config service.CircuitBreakerConfig) func(gobreaker.Counts) bool {
	return func(counts gobreaker.Counts) bool {
		if float32(counts.TotalFailures)/float32(counts.Requests) > config.Factor {
			return true
		}
		return false
	}
}

func onStateChange(name string, from gobreaker.State, to gobreaker.State) {
	fmt.Printf("%s status changed | %d -> %d", name, from, to)
}
