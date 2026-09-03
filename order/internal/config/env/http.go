package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type httpEnvConfig struct {
	ServerPort        string        `env:"HTTP_SERVER_PORT,required"`
	ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT,required"`
	RequestTimeout    time.Duration `env:"HTTP_REQUEST_TIMEOUT,required"`
	ShutdownTimeout   time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT,required"`
}

type HttpConfig struct {
	raw httpEnvConfig
}

func NewHttpConfig() (*HttpConfig, error) {
	var raw httpEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}
	return &HttpConfig{raw: raw}, nil
}

func (hc *HttpConfig) ServerPort() string {
	return hc.raw.ServerPort
}

func (hc *HttpConfig) ReadHeaderTimeout() time.Duration {
	return hc.raw.ReadHeaderTimeout
}

func (hc *HttpConfig) RequestTimeout() time.Duration {
	return hc.raw.RequestTimeout
}

func (hc *HttpConfig) ShutdownTimeout() time.Duration {
	return hc.raw.ShutdownTimeout
}
