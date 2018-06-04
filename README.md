# hyperion

Experimental Go library for abstracting out "deployment engines" like Marathon, k8s, etc.

## Example usage

```go
// appDeployerConfig := AppDeployerConfig{
// 	Type:    "marathon",
// 	Address: "http://127.0.0.1:8080",
// }
appDeployerConfig := AppDeployerConfig{
	Type:    "kubernetes",
	Address: "kubeconfig",
}
// appDeployerConfig := AppDeployerConfig{
// 	Type:    "dockerswarm",
// 	Address: "http://127.0.0.1:2377",
// }
// appDeployerConfig := AppDeployerConfig{
// 	Type:    "nomad",
// 	Address: "http://127.0.0.1:4646",
// }

appDeployer, err := NewAppDeployer(appDeployerConfig)
if err != nil {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	os.Exit(1)
}
app := GetApp(appID)
operation, err := appDeployer.DeployApp(app)
if err != nil {
	fmt.Fprintf(os.Stderr, "DeployApp error: %s\n", err)
}
fmt.Printf("operation = %+v\n", operation)
err = WaitForCompletion(ctx, operation)
if err != nil {
	fmt.Fprintf(os.Stderr, "WaitForCompletion error: %s\n", err)
}
```
