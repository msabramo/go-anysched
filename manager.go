package hyperion

import (
	"fmt"

	"git.corp.adobe.com/abramowi/hyperion/core"
)

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

type newManagerFuncType func(managerAddress string) (Manager, error)

var gRegistry = make(map[string]newManagerFuncType)

// RegisterManagerType registers the name given by managerType with a NewManager function.
func RegisterManagerType(managerType string, f newManagerFuncType) {
	if _, alreadyExists := gRegistry[managerType]; alreadyExists {
		panic(fmt.Sprintf("hyperion.RegisterManagerType: %q is already registered!", managerType))
	}
	gRegistry[managerType] = f
	ManagerTypes = append(ManagerTypes, managerType)
}

// NewManager takes a ManagerConfig and returns a specific type of Manager for
// the scheduler that the user requested (e.g.: Kubernetes, Marathon, etc.).
func NewManager(managerConfig ManagerConfig) (manager Manager, err error) {
	newManagerFunc, ok := gRegistry[managerConfig.Type]
	if !ok {
		return nil, unknownAppManagerTypeError(managerConfig.Type)
	}
	return newManagerFunc(managerConfig.Address)
}

// ManagerConfig is a struct containing configuration info that a user passes to
// the NewManager function.
type ManagerConfig struct {
	Type    string // e.g.: "marathon", "kubernetes", etc.
	Address string // e.g.: "http://127.0.0.1:8080"
}

// ManagerTypes is a slice with valid manager type names.
var ManagerTypes = []string{}

// Operation is an interface that abstracts operations executed by a Manager,
// such as deploying or destroying a service in a scheduler.
//
// Operation has methods that allow client code to check the operation's status
// or wait for it to complete.
//
// Many methods of Manager will return an Operation.
type Operation = core.Operation

// Svc is short for "service" and it is our term for something that gets
// scheduled or destroyed by a Manager
// e.g.: a Marathon application, Kubernetes deployment, etc.
// Svc is our abstract term that encompasses these types of things

// SvcCfg is used to pass information to a Manager about how to configure a
// service.
//
// In other words, it serves as an input to various Manager methods.
type SvcCfg = core.SvcCfg

// Svc contains information about a service, such as when it was started and
// how many tasks are running.
type Svc = core.Svc

// Task contains information about an individual task, such as when it was
// started and what IP addresses are assigned to it.
type Task = core.Task

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
