package hyperion

type newManagerFuncType func(managerAddress string) (Manager, error)

// gManagerTypeRegistry is a map of manager type names to functions that create
// new managers
var gManagerTypeRegistry = make(map[string]newManagerFuncType)

// gManagerTypes is a slice with valid manager type names.
var gManagerTypes = []string{}

// ClearManagerTypeRegistry clears the manager type registry, which is probably
// only useful for tests.
func ClearManagerTypeRegistry() {
	gManagerTypeRegistry = make(map[string]newManagerFuncType)
}

// RegisterManagerType registers the name given by managerType with a NewManager function.
func RegisterManagerType(managerType string, f newManagerFuncType) {
	if _, alreadyExists := gManagerTypeRegistry[managerType]; alreadyExists {
		panic(appManagerTypeAlreadyRegisteredError(managerType))
	}
	gManagerTypeRegistry[managerType] = f
	gManagerTypes = append(gManagerTypes, managerType)
}

// NewManager takes a ManagerConfig and returns a specific type of Manager for
// the scheduler that the user requested (e.g.: Kubernetes, Marathon, etc.).
func NewManager(managerConfig ManagerConfig) (manager Manager, err error) {
	newManagerFunc, ok := gManagerTypeRegistry[managerConfig.Type]
	if !ok {
		return nil, appManagerTypeUnknownError(managerConfig.Type)
	}
	return newManagerFunc(managerConfig.Address)
}
