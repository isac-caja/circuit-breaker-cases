package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/isac-caja/circuit-breaker-cases/app/service"
)

type Handler interface {
	HandleProcess(http.ResponseWriter, *http.Request)
	HandlePing(http.ResponseWriter, *http.Request)
}

type ServiceHandler struct {
	*service.Service
}

func NewHandler(service *service.Service) Handler {
	return &ServiceHandler{
		Service: service,
	}
}

func (h *ServiceHandler) HandlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "pong")
}

func (h *ServiceHandler) HandleProcess(w http.ResponseWriter, r *http.Request) {
	ctx, buff := withLogger(r.Context())

	status, err := h.Process(ctx)

	if len(buff.Bytes()) != 0 {
		defer fmt.Print(buff)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(status)
		io.WriteString(w, "ok")
	}
}

func withLogger(ctx context.Context) (context.Context, *bytes.Buffer) {
	buff := new(bytes.Buffer)
	logger := log.New(buff, "Process: ", log.Lshortfile)
	ctx = context.WithValue(ctx, "logger", logger)
	ctx = context.WithValue(ctx, "buffLogger", buff)
	return ctx, buff
}
