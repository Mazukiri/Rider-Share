/*
Package retry provides a simple retry mechanism with exponential backoff.
It is as abstract as possible to allow for different retry strategies.
*/
package retry

import (
	"context"
	"log"
	"time"
)

type Config struct {
	MaxRetries  int
	InitialWait time.Duration
	MaxWait     time.Duration
}

// DefaultConfig returns a Config with sensible default values
func DefaultConfig() Config {
	return Config{
		MaxRetries:  3,
		InitialWait: 1 * time.Second,
		MaxWait:     10 * time.Second,
	}
}