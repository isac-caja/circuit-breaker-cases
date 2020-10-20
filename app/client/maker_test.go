package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/isac-caja/circuit-breaker-cases/app/client"
	"github.com/isac-caja/circuit-breaker-cases/app/service"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestClientMaker(t *testing.T) {
	t.Run("Should create Clients by config", func(t *testing.T) {
		ctx := context.Background()
		config := service.ClientConfigList{
			{Endpoint: "http://alpha.com/beta", Timeout: 100},
			{Endpoint: "http://charlie.com/delta", Timeout: 200},
		}

		defer gock.Off()
		gock.
			New("http://alpha.com/").
			Get("beta").
			Reply(200)

		gock.
			New("http://charlie.com/").
			Get("delta").
			Reply(200)

		maker := client.NewMaker(config)
		clients := maker.Make()
		assert.NotNil(t, clients)

		for _, c := range clients {
			status, err := c.Get(ctx)
			assert.Equal(t, http.StatusOK, status)
			assert.Nil(t, err)
		}
	})
}
