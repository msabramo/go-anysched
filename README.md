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

## CLI

This repo also comes with a CLI that allows you to exercise some of the
library's features.

Run `make build` and a binary will be built called `hyperion-cli`. You can run it to get help:

```
$ make build
go build ./cmd/hyperion-cli
CLICOLOR=1 ls -l hyperion-cli
-rwxr-xr-x  1 abramowi  staff  44578428 Jul 17 15:53 hyperion-cli

$ ./hyperion-cli
A command that demos the hyperion library, allowing the user
to deploy services to Marathon, Kubernetes, etc.

Usage:
  hyperion-cli [command]

Available Commands:
  help        Help about any command
  svc         Commands for managing services
  task        Commands for managing tasks

Flags:
      --config string   config file (default is $HOME/.hyperion-cli.yaml)
  -e, --env string      environment to target
  -h, --help            help for hyperion-cli

Use "hyperion-cli [command] --help" for more information about a command.
```

Some things you can do:

### Deploy a service

```
./hyperion-cli svc deploy --svc-id=httpbin --image=citizenstig/httpbin:latest --count=3
```

### Destroy a service

```
./hyperion-cli svc destroy --svc-id=httpbin
```


[examples]: examples
