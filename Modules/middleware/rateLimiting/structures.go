package ratelimiting

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	mu      sync.Mutex
	clients = make(map[string]*Client)
)

type Client struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}
