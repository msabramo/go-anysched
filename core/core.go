package core

import (
	"context"
	"time"
)

type App struct {
	ID    string
	Image string
	Count int
}

type Operation interface{}

type AsyncOperation interface {
	Wait(ctx context.Context, timeout time.Duration) error // Wait for async operation to finish and return error or nil
}
