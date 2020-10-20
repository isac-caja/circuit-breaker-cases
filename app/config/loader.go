package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/isac-caja/circuit-breaker-cases/app/service"
)

type (
	Loader interface {
		Load() (*service.Config, error)
	}

	YamlConfig struct {
		RequestMode string `yaml:"request_mode"`
		Service     struct {
			WaitTime     int `yaml:"wait_time"`
			ReturnStatus int `yaml:"return_status"`
		}
		Requests []struct {
			ServiceName        string `yaml:"service_name"`
			Path               string
			TimeoutLimit       int     `yaml:"timeout_limit"`
			CircuitBreakerType string  `yaml:"circuit_breaker_type"`
			MaxRequests        int     `yaml:"circuit_breaker_max_requests"`
			Interval           int     `yaml:"circuit_breaker_interval"`
			Timeout            int     `yaml:"circuit_breaker_timeout"`
			Factor             float32 `yaml:"circuit_breaker_factor"`
			AllowFailures      int     `yaml:"circuit_breaker_allow_failures"`
		}
	}
)

type YamlLoader struct {
	filePath string
}

func New() (*service.Config, error) {
	f := os.Getenv("CONFIG_FILE_PATH")
	if f == "" {
		f = "/Users/iscajades/Projects/circuit-breaker-cases/resources/config.test.yaml"
	}
	loader := NewYamlLoader(f)
	return loader.Load()
}

func NewYamlLoader(filePath string) Loader {
	return &YamlLoader{
		filePath: filePath,
	}
}

func (l *YamlLoader) Load() (*service.Config, error) {
	absPath, err := filepath.Abs(l.filePath)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	yamlConfig := new(YamlConfig)
	if e := yaml.Unmarshal(data, yamlConfig); e != nil {
		return nil, e
	}

	config := new(service.Config)
	config.RequesterConfig.RequesterMode = parseMode(yamlConfig.RequestMode)
	config.WaitTime = time.Duration(yamlConfig.Service.WaitTime) * time.Millisecond
	config.ReturnStatus = yamlConfig.Service.ReturnStatus
	for _, req := range yamlConfig.Requests {
		var cb service.CircuitBreakerConfig
		if req.CircuitBreakerType != "" {
			cb = service.CircuitBreakerConfig{
				CircuitBreakerStrategy: parseStrategy(req.CircuitBreakerType),
				MaxRequests:            uint32(req.MaxRequests),
				Interval:               time.Duration(req.Interval) * time.Second,
				Timeout:                time.Duration(req.Timeout) * time.Second,
				Factor:                 req.Factor,
				AllowFailures:          uint32(req.AllowFailures),
			}
		}
		cl := service.ClientConfig{
			Endpoint:             fmt.Sprintf("http://%s/%s", req.ServiceName, req.Path),
			Timeout:              time.Duration(req.TimeoutLimit),
			CircuitBreakerConfig: cb,
		}
		config.ClientConfigList = append(config.ClientConfigList, cl)
	}
	return config, nil
}

func parseStrategy(name string) service.CircuitBreakerStrategy {
	switch name {
	case "fixed":
		return service.Fixed
	case "percentage":
		return service.Percentage
	default:
		panic(fmt.Errorf("Invalid circuit breaker strategy type: %s", name))
	}
}

func parseMode(name string) service.RequesterMode {
	switch name {
	case "async":
		return service.Async
	default:
		return service.Sync
	}
}
