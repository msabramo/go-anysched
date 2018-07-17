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

Run `make build` and a binary will be built called `bin/hyperion-cli`:

```
$ make build
go build -o bin/hyperion-cli ./cmd/hyperion-cli
CLICOLOR=1 ls -l bin/hyperion-cli
-rwxr-xr-x  1 abramowi  staff  44578428 Jul 17 16:26 bin/hyperion-cli
```

You can run it to get help:

```
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

If you have Marathon running, then you can run `make cli-smoketest-marathon`:

```
$ make cli-smoketest-marathon
...
--------------------------------------------------------------------------------
Deploying service in local_marathon ...
--------------------------------------------------------------------------------

bin/hyperion-cli svc deploy --svc-id=hyperion-cli-test-20180717164354 --image=k8s.gcr.io/echoserver:1.4 --count=1
Using config file: /Users/abramowi/go/src/git.corp.adobe.com/abramowi/hyperion/hyperion-cli.yaml
marathonDeploymentID           : 173497ab-9d3d-4fe1-8b9a-e0e6377219fe

[2018-07-17T16:43:55-07:00] Not all tasks running. 0 task(s) running.
[2018-07-17T16:43:56-07:00] Not all tasks running. 0 task(s) running.
[2018-07-17T16:43:57-07:00] All tasks running. 1 task(s) running.
Deployment completed in 3.356941183s

- name: hyperion-cli-test-20180717164354.4aa0eb1e-8a1b-11e8-9076-0242ac120005
  app-id: /hyperion-cli-test-20180717164354
  host-name: slave2.192.168.65.3.xip.io
  ip-addresses:
  - 172.17.0.2
  ports:
  - 12021
  mesos-slave-id: bfb73df5-daf4-49aa-8ed1-af4b7a259ff9-S0
  stage-time: 2018-07-17T23:44:05.562Z
  start-time: 2018-07-17T23:44:07.979Z
  state: TASK_RUNNING
  version: 2018-07-17T23:44:05.368Z

--------------------------------------------------------------------------------
Destroying service in local_marathon ...
--------------------------------------------------------------------------------

bin/hyperion-cli svc destroy --svc-id=hyperion-cli-test-20180717164354
Using config file: /Users/abramowi/go/src/git.corp.adobe.com/abramowi/hyperion/hyperion-cli.yaml
Service "hyperion-cli-test-20180717164354" deleted.
```

If you have Kubernetes running (e.g.: [minikube]), then you can run `make
cli-smoketest-kubernetes`:

```
$ make cli-smoketest-kubernetes
HYPERIONCLI_ENV=kubeconfig /Library/Developer/CommandLineTools/usr/bin/make _cli-smoketest

--------------------------------------------------------------------------------
Deploying service in kubeconfig ...
--------------------------------------------------------------------------------

bin/hyperion-cli svc deploy --svc-id=hyperion-cli-test-20180717164028 --image=k8s.gcr.io/echoserver:1.4 --count=1
Using config file: /Users/abramowi/go/src/git.corp.adobe.com/abramowi/hyperion/hyperion-cli.yaml
name                           : hyperion-cli-test-20180717164028
creationTimestamp              : 2018-07-17T15:53:39-07:00
resourceVersion                : 545751
selfLink                       : /apis/apps/v1/namespaces/default/deployments/hyperion-cli-test-20180717164028
spec.strategy                  : {RollingUpdate &RollingUpdateDeployment{MaxUnavailable:25%,MaxSurge:25%,}}
uid                            : 3f2e6af8-8a14-11e8-ba27-08002786bb43
namespace                      : default
generation                     : 1

[2018-07-17T15:53:39-07:00] Waiting for deployment "hyperion-cli-test-20180717164028" to finish: 0 of 1 updated replicas are available...
[2018-07-17T15:53:41-07:00] Deployment "hyperion-cli-test-20180717164028" successfully rolled out. 1 of 1 updated replicas are available.
Deployment completed in 2.108763226s

- name: hyperion-cli-test-20180717164028-74d75f6d56-jlf55
  host-ip: 10.0.2.15
  task-ip: 172.17.0.10
  ready-time: 2018-07-17T15:53:41-07:00

--------------------------------------------------------------------------------
Destroying service in kubeconfig ...
--------------------------------------------------------------------------------

bin/hyperion-cli svc destroy --svc-id=hyperion-cli-test-20180717164028
Using config file: /Users/abramowi/go/src/git.corp.adobe.com/abramowi/hyperion/hyperion-cli.yaml
Service "hyperion-cli-test-20180717164028" deleted.
```

If you have both Marathon and Kubernetes running, then you can run `make
cli-smoketest`. This will run the above tests for both Marathon and Kubernetes.


[examples]: examples
[minikube]: https://github.com/kubernetes/minikube
