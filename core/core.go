package core

import (
	"context"
	"time"
)

type SvcCfg struct {
	ID    string
	Image string
	Count int

	DeployTimeoutDuration *time.Duration // pointer because optional
}

type Svc struct {
	ID             string     `yaml:"ID" json:"ID"`
	TasksRunning   *int       `yaml:"tasks-running,omitempty" json:"tasks-running,omitempty"`
	TasksHealthy   *int       `yaml:"tasks-healthy,omitempty" json:"tasks-healthy,omitempty"`
	TasksUnhealthy *int       `yaml:"tasks-unhealthy,omitempty" json:"tasks-unhealthy,omitempty"`
	CreationTime   *time.Time `yaml:"creation-time,omitempty" json:"creation-time,omitempty"`
}

type OperationStatus struct {
	ClientTime         time.Time
	LastTransitionTime time.Time
	LastUpdateTime     time.Time
	Msg                string
	Done               bool
}

type Task struct {
	Name                string     `yaml:"name" json:"name"`
	AppID               string     `yaml:"app-id,omitempty" json:"app-id,omitempty"`
	HostName            string     `yaml:"host-name,omitempty" json:"host-name,omitempty"`
	HostIP              string     `yaml:"host-ip,omitempty" json:"host,omitempty"`
	TaskIP              string     `yaml:"task-ip,omitempty" json:"task-ip,omitempty"`
	IPAddresses         []string   `yaml:"ip-addresses,omitempty" json:"ip-addresses,omitempty"`
	Ports               []int      `yaml:"ports,omitempty" json:"ports,omitempty"`
	ServicePorts        []int      `yaml:"service-ports,omitempty" json:"service-ports,omitempty"`
	MesosSlaveID        string     `yaml:"mesos-slave-id,omitempty" json:"mesos-slave-id,omitempty"`
	StageTime           *time.Time `yaml:"stage-time,omitempty" json:"stage-time,omitempty"`
	StartTime           *time.Time `yaml:"start-time,omitempty" json:"start-time,omitempty"`
	ReadyTime           *time.Time `yaml:"ready-time,omitempty" json:"ready-time,omitempty"`
	LastHealthCheckTime *time.Time `yaml:"last-health-check-time,omitempty" json:"last-health-check-time,omitempty"`
	LastHealthyTime     *time.Time `yaml:"last-healthy-time,omitempty" json:"last-healthy-time,omitempty"`
	State               string     `yaml:"state,omitempty" json:"state,omitempty"`
	Version             string     `yaml:"version,omitempty" json:"version,omitempty"`
}

type Operation interface {
	// GetProperties returns a map with all labels, annotations, and basic
	// properties like name or uid
	GetProperties() map[string]interface{}

	// Wait waits for an operation to finish and return error or nil
	Wait(ctx context.Context) (result interface{}, err error)

	// GetStatus is for polling the status of the deployment
	GetStatus() (status *OperationStatus, err error)
}
