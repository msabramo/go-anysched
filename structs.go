package anysched

import (
	"time"
)

// ManagerConfig is a struct containing configuration info that a user passes to
// the NewManager function.
type ManagerConfig struct {
	Type    string // e.g.: "marathon", "kubernetes", etc.
	Address string // e.g.: "http://127.0.0.1:8080"
}

// Svc is short for "service" and it is our term for something that gets
// scheduled or destroyed by a Manager
// e.g.: a Marathon application, Kubernetes deployment, etc.
// Svc is our abstract term that encompasses these types of things

// SvcCfg is used to pass information to a Manager about how to configure a
// service.
//
// In other words, it serves as an input to various Manager methods.
type SvcCfg struct {
	ID    string
	Image string
	Count int

	DeployTimeoutDuration *time.Duration // pointer because optional
}

// Svc contains information about a service, such as when it was started and
// how many tasks are running.
type Svc struct {
	ID             string     `yaml:"ID" json:"ID"`
	TasksRunning   *int       `yaml:"tasks-running,omitempty" json:"tasks-running,omitempty"`
	TasksHealthy   *int       `yaml:"tasks-healthy,omitempty" json:"tasks-healthy,omitempty"`
	TasksUnhealthy *int       `yaml:"tasks-unhealthy,omitempty" json:"tasks-unhealthy,omitempty"`
	CreationTime   *time.Time `yaml:"creation-time,omitempty" json:"creation-time,omitempty"`
}

// OperationStatus represents the status of a pending operation, such as a deployment.
type OperationStatus struct {
	ClientTime         time.Time
	LastTransitionTime time.Time
	LastUpdateTime     time.Time
	Msg                string
	Done               bool
}

// Task contains information about an individual task, such as when it was
// started and what IP addresses are assigned to it.
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
