package client

import (
	"context"
	"net/http"

	"github.com/isac-caja/circuit-breaker-cases/app/service"
)

func New(config service.ClientConfig) service.Client {
	client := &http.Client{Timeout: config.Timeout}
	return &Client{
		Client:       client,
		ClientConfig: config,
	}
}

func NewMaker(config service.ClientConfigList) *ClientMaker {
	return &ClientMaker{
		ClientConfigList: config,
	}
}

type ClientMaker struct {
	service.ClientConfigList
}

func (m *ClientMaker) Make() service.Clients {
	var clients service.Clients
	for _, conf := range m.ClientConfigList {
		client := New(conf)
		breaker := NewBreaker(conf.CircuitBreakerConfig)
		clients = append(clients, WithBreaker(client, breaker))
	}
	return clients
}

type Client struct {
	*http.Client
	service.ClientConfig
}

func (c *Client) Get(ctx context.Context) (int, error) {
	resp, err := c.Client.Get(c.Endpoint)
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}
