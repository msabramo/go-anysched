# anysched

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

Run `make build` and a binary will be built called `bin/anysched-cli`:

```
$ make build
go build -o bin/anysched-cli ./cmd/anysched-cli
CLICOLOR=1 ls -l bin/anysched-cli
-rwxr-xr-x  1 abramowi  staff  44578428 Jul 17 16:26 bin/anysched-cli
```

You can run it to get help:

```
$ bin/anysched-cli
A command that demos the anysched library, allowing the user
to deploy services to Marathon, Kubernetes, etc.

Usage:
  anysched-cli [command]

Available Commands:
  help        Help about any command
  svc         Commands for managing services
  task        Commands for managing tasks

Flags:
      --config string   config file (default is $HOME/.anysched-cli.yaml)
  -e, --env string      environment to target
  -h, --help            help for anysched-cli

Use "anysched-cli [command] --help" for more information about a command.
```

Some things you can do:

### Deploy a service

```
bin/anysched-cli svc deploy --svc-id=httpbin --image=citizenstig/httpbin:latest --count=3
```

### Destroy a service

```
bin/anysched-cli svc destroy --svc-id=httpbin
```

## Unit tests

Run `make test`.

```
$ make test
go test ./...
ok  	github.com/msabramo/go-anysched	0.037s
?   	github.com/msabramo/go-anysched/cmd/anysched-cli	[no test files]
?   	github.com/msabramo/go-anysched/examples	[no test files]
ok  	github.com/msabramo/go-anysched/managers/dockerswarm	0.042s
ok  	github.com/msabramo/go-anysched/managers/kubernetes	0.070s
ok  	github.com/msabramo/go-anysched/managers/marathon	0.037s
ok  	github.com/msabramo/go-anysched/managers/nomad	0.037s
ok  	github.com/msabramo/go-anysched/utils	0.075s
```

## Unit test coverage

```
$ make test-cover
ANYSCHEDCLI_ENV=minikube scripts/coverage
ok      .                                          0.061s coverage: 100.0% of statements
ok      ./managers/dockerswarm                     0.072s coverage:  20.0% of statements
ok      ./managers/kubernetes                      8.205s coverage:  91.6% of statements
ok      ./managers/marathon                        0.061s coverage:   8.1% of statements
ok      ./managers/nomad                           0.071s coverage:  27.3% of statements
ok      ./utils                                    0.079s coverage: 100.0% of statements

real	0m11.607s
user	0m7.439s
sys	0m3.060s
Total code coverage: 58.6%
```

## Integration tests

If you have Marathon running, then you can run `make cli-smoketest-marathon`:

```
$ make cli-smoketest-marathon
...
--------------------------------------------------------------------------------
Deploying service in local_marathon ...
--------------------------------------------------------------------------------

bin/anysched-cli svc deploy --svc-id=anysched-cli-test-20180717164354 --image=k8s.gcr.io/echoserver:1.4 --count=1
Using config file: /Users/abramowi/go/src/github.com/msabramo/go-anysched/anysched-cli.yaml
marathonDeploymentID           : 173497ab-9d3d-4fe1-8b9a-e0e6377219fe

[2018-07-17T16:43:55-07:00] Not all tasks running. 0 task(s) running.
[2018-07-17T16:43:56-07:00] Not all tasks running. 0 task(s) running.
[2018-07-17T16:43:57-07:00] All tasks running. 1 task(s) running.
Deployment completed in 3.356941183s

- name: anysched-cli-test-20180717164354.4aa0eb1e-8a1b-11e8-9076-0242ac120005
  app-id: /anysched-cli-test-20180717164354
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

bin/anysched-cli svc destroy --svc-id=anysched-cli-test-20180717164354
Using config file: /Users/abramowi/go/src/github.com/msabramo/go-anysched/anysched-cli.yaml
Service "anysched-cli-test-20180717164354" deleted.
```

If you have Kubernetes running (e.g.: [minikube]), then you can run `make
cli-smoketest-kubernetes`:

```
$ make cli-smoketest-kubernetes
ANYSCHEDCLI_ENV=kubeconfig /Library/Developer/CommandLineTools/usr/bin/make _cli-smoketest

--------------------------------------------------------------------------------
Deploying service in kubeconfig ...
--------------------------------------------------------------------------------

bin/anysched-cli svc deploy --svc-id=anysched-cli-test-20180717164028 --image=k8s.gcr.io/echoserver:1.4 --count=1
Using config file: /Users/abramowi/go/src/github.com/msabramo/go-anysched/anysched-cli.yaml
name                           : anysched-cli-test-20180717164028
creationTimestamp              : 2018-07-17T15:53:39-07:00
resourceVersion                : 545751
selfLink                       : /apis/apps/v1/namespaces/default/deployments/anysched-cli-test-20180717164028
spec.strategy                  : {RollingUpdate &RollingUpdateDeployment{MaxUnavailable:25%,MaxSurge:25%,}}
uid                            : 3f2e6af8-8a14-11e8-ba27-08002786bb43
namespace                      : default
generation                     : 1

[2018-07-17T15:53:39-07:00] Waiting for deployment "anysched-cli-test-20180717164028" to finish: 0 of 1 updated replicas are available...
[2018-07-17T15:53:41-07:00] Deployment "anysched-cli-test-20180717164028" successfully rolled out. 1 of 1 updated replicas are available.
Deployment completed in 2.108763226s

- name: anysched-cli-test-20180717164028-74d75f6d56-jlf55
  host-ip: 10.0.2.15
  task-ip: 172.17.0.10
  ready-time: 2018-07-17T15:53:41-07:00

--------------------------------------------------------------------------------
Destroying service in kubeconfig ...
--------------------------------------------------------------------------------

bin/anysched-cli svc destroy --svc-id=anysched-cli-test-20180717164028
Using config file: /Users/abramowi/go/src/github.com/msabramo/go-anysched/anysched-cli.yaml
Service "anysched-cli-test-20180717164028" deleted.
```

If you have both Marathon and Kubernetes running, then you can run `make
cli-smoketest`. This will run the above tests for both Marathon and Kubernetes.

## More make targets

The `Makefile` is self-documenting. Running `make` with no target will print a list of targets:

```
$ make
build                          Build all the things
check                          Run tests and linters
clean                          Clean up files that aren't checked into version control
cli-smoketest-kubernetes       Quickly exercise anysched-cli for Marathon
cli-smoketest-marathon         Quickly exercise anysched-cli for Marathon
cli-smoketest                  Quickly exercise anysched-cli for Marathon and Kubernetes
lint                           Run golint linter
metalinter                     Run gometalinter, which does a bunch of checks
test-cover-html                Generate HTML test coverage report
test-cover                     Generate test coverage report
test-race                      Run tests with race detector
test                           Run tests
top-cyclo                      Display function with most cyclomatic complexity
vet                            Run go vet linter
```

## CLI Docker image

If you don't have Go installed or have trouble building from source, you can
try using Docker. There is a script called
[`scripts/cli-docker`](scripts/cli-docker) that will automatically build the
Docker image if you don't already have it and then run
a container using it. Example:

```
$ scripts/cli-docker svc deploy --svc-id=httpbin --image=citizenstig/httpbin:latest --count=3
Building docker image...
Sending build context to Docker daemon  360.4kB
Step 1/11 : FROM golang:alpine AS build-env
 ---> 34d3217973fd
Step 2/11 : RUN apk add --update git
 ---> Running in f2fa64a62f5b
...
Using config file: /workdir/anysched-cli.yaml
name                           : httpbin
selfLink                       : /apis/apps/v1/namespaces/default/deployments/httpbin
spec.strategy                  : {RollingUpdate &RollingUpdateDeployment{MaxUnavailable:25%,MaxSurge:25%,}}
resourceVersion                : 30081
uid                            : 096f040b-8c4c-11e8-a0ad-080027aa669d
creationTimestamp              : 2018-07-20T18:38:03Z
namespace                      : default
generation                     : 1

[2018-07-20T18:38:03Z] Waiting for deployment "httpbin" to finish: 0 of 3 updated replicas are available...
[2018-07-20T18:38:06Z] Waiting for deployment "httpbin" to finish: 1 of 3 updated replicas are available...
[2018-07-20T18:38:07Z] Waiting for deployment "httpbin" to finish: 2 of 3 updated replicas are available...
[2018-07-20T18:38:09Z] Deployment "httpbin" successfully rolled out. 3 of 3 updated replicas are available.
Deployment completed in 6.323958015s

- name: httpbin-5d7c976bcd-9kjz5
  host-ip: 10.0.2.15
  task-ip: 172.17.0.4
  ready-time: 2018-07-20T18:38:05Z
- name: httpbin-5d7c976bcd-wmsvl
  host-ip: 10.0.2.15
  task-ip: 172.17.0.5
  ready-time: 2018-07-20T18:38:07Z
- name: httpbin-5d7c976bcd-xn6dx
  host-ip: 10.0.2.15
  task-ip: 172.17.0.6
  ready-time: 2018-07-20T18:38:09Z
```


[examples]: examples
[minikube]: https://github.com/kubernetes/minikube
