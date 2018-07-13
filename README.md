# hyperion

An experimental Go library that attempts to provide a common interface for
various container-oriented app management systems -- e.g.:

- Kubernetes
- Marathon

## Example usage

```go
import (
	"git.corp.adobe.com/abramowi/hyperion"
)

managerConfig := hyperion.ManagerConfig{Type: hyperion.ManagerTypeKubernetes}
// or alternatively one of the following:
//
// managerConfig := hyperion.ManagerConfig{
// 	Type:    hyperion.ManagerTypeMarathon,
// 	Address: "http://127.0.0.1:8080",
// }
// managerConfig := hyperion.ManagerConfig{
// 	Type:    hyperion.ManagerTypeDockerSwarm,
// 	Address: "http://127.0.0.1:2377",
// }
// managerConfig := hyperion.ManagerConfig{
// 	Type:    hyperion.ManagerTypeNomad
// 	Address: "http://127.0.0.1:4646",
// }

manager, err := hyperion.NewManager(managerConfig)
if err != nil {
	return err
}
svc := hyperion.Svc{
	ID:    "my-svc-id",
	Image: "citizenstig/httpbin",
	Count: 4,
}
operation, err := manager.DeploySvc(svc)
if err != nil {
	return err
}
```
