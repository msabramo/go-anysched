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

Run `make build` and a binary will be built called `bin/hyperion-cli`. You can
run it to get help:

```
$ make build
go build -o bin/hyperion-cli ./cmd/hyperion-cli
CLICOLOR=1 ls -l bin/hyperion-cli
-rwxr-xr-x  1 abramowi  staff  44578428 Jul 17 16:26 bin/hyperion-cli

$ bin/hyperion-cli
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
bin/hyperion-cli svc deploy --svc-id=httpbin --image=citizenstig/httpbin:latest --count=3
```

### Destroy a service

```
bin/hyperion-cli svc destroy --svc-id=httpbin
```

## Unit tests

Run `make test`.

```
$ make test
go test ./...
ok  	git.corp.adobe.com/abramowi/hyperion	0.049s
?   	git.corp.adobe.com/abramowi/hyperion/cmd/hyperion-cli	[no test files]
ok  	git.corp.adobe.com/abramowi/hyperion/dockerswarm	0.040s
?   	git.corp.adobe.com/abramowi/hyperion/examples	[no test files]
ok  	git.corp.adobe.com/abramowi/hyperion/kubernetes	0.048s
ok  	git.corp.adobe.com/abramowi/hyperion/marathon	0.042s
ok  	git.corp.adobe.com/abramowi/hyperion/nomad	0.041s
ok  	git.corp.adobe.com/abramowi/hyperion/utils	0.042s
```

## Integration tests

If you have both Marathon and Kubernetes running, then you can run `make
cli-smoketest`. This uses `bin/hyperion-cli` to deploy and destroy a service,
and it will do it in both Marathon and Kubernetes.


[examples]: examples
