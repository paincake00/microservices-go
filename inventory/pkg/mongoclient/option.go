package mongoclient

import "time"

type Option func(m *Client)

func OperationTimeout(timeout time.Duration) Option {
	return func(m *Client) {
		m.operationTimeout = timeout
	}
}

func MaxPoolSize(value int) Option {
	return func(m *Client) {
		m.maxPoolSize = value
	}
}

func MinPoolSize(value int) Option {
	return func(m *Client) {
		m.minPoolSize = value
	}
}

func MaxCons(value int) Option {
	return func(m *Client) {
		m.maxCons = value
	}
}

func MaxConIdleTime(value time.Duration) Option {
	return func(m *Client) {
		m.maxConIdleTime = value
	}
}

func ShutdownTimeout(value time.Duration) Option {
	return func(m *Client) {
		m.shutdownTimeout = value
	}
}
