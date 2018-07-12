package core

import (
	"context"
	"time"
)

type App struct {
	ID    string
	Image string
	Count int

	DeployTimeoutDuration *time.Duration // pointer because optional
}

type Status struct {
	ClientTime         time.Time
	LastTransitionTime time.Time
	LastUpdateTime     time.Time
	Msg                string
	Done               bool
}

type Operation interface {
	// GetProperties returns a map with all labels, annotations, and basic
	// properties like name or uid
	GetProperties() map[string]interface{}

	// Wait waits for an operation to finish and return error or nil
	Wait(ctx context.Context) (result interface{}, err error)

	// GetStatus is for polling the status of the deployment
	GetStatus() (status *Status, err error)
}
