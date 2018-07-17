# hyperion

An experimental Go library that attempts to provide a common interface for
various container-oriented app management systems -- e.g.:

- Kubernetes
- Marathon

## Example library usage

See [examples].

Run an example:

This assumes that you have a `~/.kube/config` that points to a running
Kubernetes cluster. It will create a deployment called `my-svc-id` with 4
running pods.

```
make -C examples run-deploy-example
```


[examples]: examples
