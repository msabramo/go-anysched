# hyperion

Experimental Go library for abstracting out "deployment engines" like Marathon, k8s, etc.

## Example usage

```go
import (
	hyperionlib "git.corp.adobe.com/abramowi/hyperion/lib"
)

managerConfig := hyperionlib.ManagerConfig{
	Type:    hyperionlib.ManagerTypeKubernetes,
	Address: "kubeconfig",
}
// or alternatively one of the following:
//
// managerConfig := hyperonlib.ManagerConfig{
// 	Type:    hyperionlib.ManagerTypeMarathon,
// 	Address: "http://127.0.0.1:8080",
// }
// managerConfig := hyperonlib.ManagerConfig{
// 	Type:    hyperionlib.ManagerTypeDockerSwarm,
// 	Address: "http://127.0.0.1:2377",
// }
// managerConfig := hyperonlib.ManagerConfig{
// 	Type:    hyperionlib.ManagerTypeNomad
// 	Address: "http://127.0.0.1:4646",
// }

manager, err := hyperionlib.NewManager(managerConfig)
if err != nil {
	return err
}
app := hyperionlib.App{
	ID:    "my-app-id",
	Image: "citizenstig/httpbin",
	Count: 4,
}
operation, err := manager.DeployApp(app)
if err != nil {
	return err
}
```
