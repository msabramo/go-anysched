package hyperion

import "context"

// Manager manages various types of schedulers, such as Kubernetes, Marathon, etc.
// It is an interface that is composed of various other more fine-grained
// interfaces, because fine-grained interfaces are awesome.
//
// You create a manager by calling NewManager, passing it a ManagerConfig.
type Manager interface {
	SvcDeployer
	SvcDestroyer
	SvcsGetter
	SvcTasksGetter
	TasksGetter
}

// SvcDeployer is an interface with a method for deploying a service.
type SvcDeployer interface {
	// DeploySvc takes a SvcCfg and deploys it, returning an Operation.
	DeploySvc(SvcCfg) (Operation, error)
}

// SvcDestroyer is an interface with a method for destroying a service.
type SvcDestroyer interface {
	// DestroySvc destroys a service.
	DestroySvc(svcID string) (Operation, error)
}

// SvcsGetter is an interface with a method for getting all running services.
type SvcsGetter interface {
	// Svcs returns info about the running services.
	Svcs() ([]Svc, error)
}

// SvcTasksGetter is an interface with a method for getting all running tasks
// for a particular service.
type SvcTasksGetter interface {
	// SvcTasks returns info about the running tasks for a service.
	SvcTasks(SvcCfg) ([]Task, error)
}

// TasksGetter is an interface with a method for getting all running tasks
// across all services.
type TasksGetter interface {
	// Tasks returns info about all running tasks.
	Tasks() ([]Task, error)
}

// Operation is an interface that abstracts operations executed by a Manager,
// such as deploying or destroying a service in a scheduler.
//
// Operation has methods that allow client code to check the operation's status
// or wait for it to complete.
//
// Many methods of Manager will return an Operation.
type Operation interface {
	// GetProperties returns a map with all labels, annotations, and basic
	// properties like name or uid
	GetProperties() map[string]interface{}

	// Wait waits for an operation to finish and return error or nil
	Wait(ctx context.Context) (result interface{}, err error)

	// GetStatus is for polling the status of the deployment
	GetStatus() (status *OperationStatus, err error)
}
