package pgclient

import "time"

type Option func(*Client)

func MaxOpenCons(cons int) Option {
	return func(p *Client) {
		p.maxOpenCons = cons
	}
}

func MinIdleCons(cons int) Option {
	return func(p *Client) {
		p.minIdleCons = cons
	}
}

func MaxConIdleTime(duration time.Duration) Option {
	return func(p *Client) {
		p.maxConIdleTime = duration
	}
}

func MaxConLifetime(duration time.Duration) Option {
	return func(p *Client) {
		p.maxConLifetime = duration
	}
}
