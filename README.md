# hyperion

Experimental Go library for abstracting out "deployment engines" like Marathon, k8s, etc.

## Example usage

```go
import (
	hyperionlib "git.corp.adobe.com/abramowi/hyperion/lib"
)

appDeployerConfig := hyperionlib.AppDeployerConfig{
	Type:    "kubernetes",
	Address: "kubeconfig",
}
// or alternatively one of the following:
//
// appDeployerConfig := hyperionlib.AppDeployerConfig{
// 	Type:    "marathon",
// 	Address: "http://127.0.0.1:8080",
// }
// appDeployerConfig := hyperionlib.AppDeployerConfig{
// 	Type:    "dockerswarm",
// 	Address: "http://127.0.0.1:2377",
// }
// appDeployerConfig := hyperionlib.AppDeployerConfig{
// 	Type:    "nomad",
// 	Address: "http://127.0.0.1:4646",
// }

appDeployer, err := hyperionlib.NewAppDeployer(appDeployerConfig)
if err != nil {
	return err
}
app := hyperionlib.App{
	ID:    "my-app-id",
	Image: "citizenstig/httpbin",
	Count: 4,
}
operation, err := appDeployer.DeployApp(app)
if err != nil {
	return err
}
```
