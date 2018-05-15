package core

import (
	"context"
	"time"
)

type Operation interface{}

type AsyncOperation interface {
	Wait(ctx context.Context, timeout time.Duration) error // Wait for async operation to finish and return error or nil
}
