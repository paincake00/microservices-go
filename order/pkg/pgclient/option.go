package pgclient

import "time"

type Option func(*PostgresClient)

func MaxOpenCons(cons int) Option {
	return func(p *PostgresClient) {
		p.maxOpenCons = cons
	}
}

func MinIdleCons(cons int) Option {
	return func(p *PostgresClient) {
		p.minIdleCons = cons
	}
}

func MaxConIdleTime(duration time.Duration) Option {
	return func(p *PostgresClient) {
		p.maxConIdleTime = duration
	}
}

func MaxConLifetime(duration time.Duration) Option {
	return func(p *PostgresClient) {
		p.maxConLifetime = duration
	}
}
