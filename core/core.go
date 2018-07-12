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

type TaskInfo struct {
	Name                string     `yaml:"name" json:"name"`
	HostIP              string     `yaml:"host-ip,omitempty" json:"host,omitempty"`
	TaskIP              string     `yaml:"task-ip,omitempty" json:"task-ip,omitempty"`
	StartTime           *time.Time `yaml:"start-time,omitempty" json:"start-time,omitempty"`
	ReadyTime           *time.Time `yaml:"ready-time,omitempty" json:"ready-time,omitempty"`
	LastHealthCheckTime *time.Time `yaml:"last-health-check-time,omitempty" json:"last-health-check-time,omitempty"`
	LastHealthyTime     *time.Time `yaml:"last-healthy-time,omitempty" json:"last-healthy-time,omitempty"`
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
