# Jobs as a Service (JaaS)

> A CLI for running jobs (ad-hoc containers/tasks) on Docker Swarm

This project provides a simple Golang CLI tool that binds to the Docker Swarm API to create an ad-hoc/one-shot Service and then poll until it exits. Service logs can also be retrieved if the Docker daemon API version is greater than 1.29 or if the experimental feature is enabled on the Docker daemon.

[![Build Status](https://travis-ci.org/alexellis/jaas.svg?branch=master)](https://travis-ci.org/alexellis/jaas)

**Motivation and context**

For a blog post covering use-cases for JaaS and more on the portions of the Docker API used see below:

* [Blog: One-shot containers on Docker Swarm](http://blog.alexellis.io/containers-on-swarm/)

**See also: Serverless**

If you would like to build Serverless applications with Docker Swarm or Kubernetes checkout my write-up on OpenFaaS:

* [OpenFaaS.com](https://www.openfaas.com)

The OpenFaaS project has dozens of contributors and thousands of GitHub stars - if you're here because you want to run short-lived functions then I highly recommend checking out OpenFaaS now.

## Contributions are welcome

See the [contributing guide](CONTRIBUTING.md) and do not raise a PR unless you've read it all.

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

```
# jaas run -r --image alexellis2/cows:latest
```

The `-rm` flag removes the Swarm service that was used to run your container.

> The exit code from your container will also be available, you can check it with `echo $?`

* Hiding logs

If you aren't interested in the output logs then run it with the `--show-logs=false` override:

```
# jaas run --image alexellis2/cows:latest --show-logs=false
```

* Override the command of the container:

```
# jaas run -r --image alpine:3.6 --command "uname -a"

Printing service logs
w2018-02-06T13:40:00.131678932Z Linux f56d298c4ab9 4.9.75-linuxkit-aufs #1 SMP Tue Jan 9 10:58:17 UTC 2018 x86_64 Linux
```

* Removing service after completion

To remove the service after it completes, run with the `--remove` or `-r` flag:

```
# jaas run --image alexellis2/href-counter:latest --env url=http://blog.alexellis.io/

Service created: peaceful_shirley (uva6bcqyubm1b4c80dghjhb44)
ID:  uva6bcqyubm1b4c80dghjhb44  Update at:  2017-03-14 22:19:54.381973142 +0000 UTC
...

Exit code: 0
State: complete

Printing service logs
?2017-03-14T22:19:55.660902727Z com.docker.swarm.node.id=b2dqydhfavwezorhkqi11f962,com.docker.swarm.service.id=uva6bcqyubm1b4c80dghjhb44,com.docker.swarm.task.id=yruxuawdipz2v5n0wvvm8ib0r {"internal":42,"external":2}

Removing service...
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

```
# echo "p4ssword" > pass_file
# docker secret create my_secret pass_file
# jaas run --image alexellis2/href-counter:latest --env url=http://blog.alexellis.io/ --secret my_secret --command "cat /run/secrets/my_secret"

Service created: priceless_tesla (f8gheat9f3b8cnnsjy9dth9y7)
ID:  f8gheat9f3b8cnnsjy9dth9y7  Update at:  2018-06-29 16:41:13.723257461 +0000 UTC
...........

Exit code: 0
State: complete


Printing service logs
(2018-06-29T16:41:19.057738088Z p4ssword

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

### Roadmap:

Here are several features / enhancements on the roadmap, please make additional suggestions through Github issues.

* [x] Optionally delete service after fetching exit code/logs
* [x] Support passing environmental variables
* [x] Support private registry auth via `-registryAuth` flag
* [x] Move to cobra flags/args package for CLI
* [x] Support constraints on where to run tasks
* [x] Bind-mounting volumes
* [x] Overriding container command

Todo:

* [ ] Support optional secrets through CLI flag
* [ ] Validation around images which are not in local library
* [ ] Extract stdout/stderr etc from logs in human readable format similar to `docker logs`

### Future:

* When task logs are available in the API this will be used instead of service logs.
* When event streams are released they will prevent the need to poll continually
