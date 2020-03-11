# Jobs as a Service (JaaS)

Run jobs (tasks/one-shot containers) on Docker Swarm

This project provides a simple Golang CLI tool that binds to the Docker Swarm API to create an ad-hoc/one-shot Service and then poll until it exits. Service logs can also be retrieved if the Docker daemon API version is greater than 1.29 or if the experimental feature is enabled on the Docker daemon.

[![Build Status](https://travis-ci.org/alexellis/jaas.svg?branch=master)](https://travis-ci.org/alexellis/jaas)

## Motivation and context

For a blog post covering use-cases for JaaS and more on the portions of the Docker API used see below:

* [Blog: One-shot containers on Docker Swarm](http://blog.alexellis.io/containers-on-swarm/)

Use-cases:

* Use an elastic cluster as your computer
* Clean up DB indexes
* Send emails
* Batch processing
* Replace cron scripts
* Run your server maintenance tasks
* Schedule dev-ops tasks

### See also: Serverless

If you would like to build Serverless applications with Kubernetes or Docker Swarm checkout OpenFaaS:

* [OpenFaaS.com](https://www.openfaas.com)

> The OpenFaaS project has dozens of contributors and thousands of GitHub stars - if you're here because you want to run short-lived functions then I highly recommend checking out OpenFaaS now.

## Get started

### Build and install the code

Pre-requisites:

* Docker 1.13 or newer (experimental mode must be enabled if accessing service logs with Docker versions >= 1.13 and < 1.29)
* [Go 1.9.2 (or Golang container)](https://golang.org/dl/)
* Enable Swarm Mode (`docker swarm init`)

**Run these commands**

```
# export GOPATH=$HOME/go
# go get -d -v github.com/alexellis/jaas
# cd $GOPATH/src/github.com/alexellis/jaas
# go install
# export PATH=$PATH:$GOPATH/bin
```

Now test `jaas` with `jaas --help`

### Running a task / batch job / one-shot container

* Run your first one-shot container with `jaas run`:

```bash
# jaas run -r --image alexellis2/cows:latest
```

The `-r` flag removes the Swarm service that was used to run your container.

> The exit code from your container will also be available, you can check it with `echo $?`

* Hiding logs

If you aren't interested in the output logs then run it with the `--show-logs=false` override:

```bash
# jaas run --image alexellis2/cows:latest --show-logs=false
```

* Override the command of the container

```bash
# jaas run --image alpine:3.8 --command "uname -a"

Printing service logs
w2018-02-06T13:40:00.131678932Z Linux f56d298c4ab9 4.9.75-linuxkit-aufs #1 SMP Tue Jan 9 10:58:17 UTC 2018 x86_64 Linux
```

You can also try the example in `examples/gotask`:


```bash
# jaas run -r --image alexellis2/go-task:2020-03-11
```

* Environment variables

Set environment variables with `--env` or `-e`:

```bash
# jaas run --image alpine:3.8 --env ENV1=val1 --env ENV2=val2 --command "env"

Service created: inspiring_elion (j90qjtc14usgps9t60tvogmts)
ID:  j90qjtc14usgps9t60tvogmts  Update at:  2018-07-14 18:02:57.147797437 +0000 UTC
...........

Exit code: 0
State: complete


Printing service logs
a2018-07-14T18:03:01.465983797Z PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
52018-07-14T18:03:01.466098037Z HOSTNAME=de0b5614fc88
)2018-07-14T18:03:01.466111965Z ENV1=val1
)2018-07-14T18:03:01.466122558Z ENV2=val2
*2018-07-14T18:03:01.466132520Z HOME=/root

Removing service...
```

* Removing service after completion

By default, the service is removed after it completes. To prevent that, run with the `--remove` or `-r` flag set to `false`:

```bash
# jaas run --image alpine:3.8 --remove=false

Service created: zen_hoover (nwf2zey3i387zkx5gp7yjk053)
ID:  nwf2zey3i387zkx5gp7yjk053  Update at:  2018-07-08 20:19:39.320494122 +0000 UTC
............

Exit code: 0
State: complete


Printing service logs

# docker service ls

ID            NAME        MODE        REPLICAS  IMAGE       PORTS
nwf2zey3i387  zen_hoover  replicated  0/1       alpine:3.7
```

* Docker authentication for registries

You can use `jaas` with Docker images in private registries or registries which require authentication.

Just run `docker login` then pass the `--registry` parameter and the encoded string you find in `~/.docker/config.json`.

If you want to encode a string manually then do the following:

```bash
$ export auth='{
    "username" : "myUserName",
    "password" : "secret",
    "email" : "my@email",
    "serveraddress" : "my.reg.domain"
  }'
$ jaas run --registry="`echo $auth | base64`" --image my.reg.domain/hello-world:latest
```

* Adding secret to service

To give the service access to an _existing secret_. run with the `--secret` or `-s` flag:

```bash
$ echo -n "S3_ACCESS_KEY_HERE" | docker secret create s3-access-key -
$ jaas run --image alpine:3.7 --secret s3-access-key --command "cat /run/secrets/s3-access-key"

Service created: priceless_tesla (f8gheat9f3b8cnnsjy9dth9y7)
ID:  f8gheat9f3b8cnnsjy9dth9y7  Update at:  2018-06-29 16:41:13.723257461 +0000 UTC
...........

Exit code: 0
State: complete

Printing service logs
(2018-06-29T16:41:19.057738088Z S3_ACCESS_KEY_HERE

Removing service...
```

_Notes on images_

You can have a multi-node swarm but make sure whatever image you choose is available in an accessible registry.

> A local image will not need to be pushed to a registry.

* Running jaas in a container

You can also run `jaas` in a container, but the syntax becomes slightly more verbose:

```
# docker run -ti -v /var/run/docker.sock:/var/run/docker.sock \
  alexellis2/jaas run --image alexellis2/cows:latest
```

## Real-life example

You can use `jaas` to get the value of your OpenFaaS gateway password on Docker Swarm. See the [OpenFaaS troubleshooting guide](https://docs.openfaas.com/deployment/troubleshooting/#swarm_1) for the usage.

### Roadmap:

Here are several features / enhancements on the roadmap, please make additional suggestions through Github issues.

* [x] Optionally delete service after fetching exit code/logs
* [x] Support passing environmental variables
* [x] Support private registry auth via `-registryAuth` flag
* [x] Move to cobra flags/args package for CLI
* [x] Support constraints on where to run tasks
* [x] Bind-mounting volumes
* [x] Overriding container command
* [x] Support optional secrets through CLI flag

Todo:

* [ ] Validation around images which are not in local library
* [ ] Extract stdout/stderr etc from logs in human readable format similar to `docker logs`
* [ ] Support incoming [jobs API for Swarm](https://github.com/moby/moby/issues/39447)

### Future:

* When task logs are available in the API this will be used instead of service logs.
* When event streams are released they will prevent the need to poll continually

## Similar tools

* [kjob](https://github.com/stefanprodan/kjob) by Stefan Prodan appears to be a close variation of jaas, but for Kubernetes.

## Contributions are welcome

See the [contributing guide](CONTRIBUTING.md) and do not raise a PR unless you've read it all.
