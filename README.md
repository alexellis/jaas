# Ad-hoc Jobs as a Service (JaaS)

This project provides a simple Golang CLI tool that binds to the Docker Swarm API to create an ad-hoc/one-shot Service and then poll until it exits. Service logs can also be retrieved if the experimental feature is enabled on the Docker daemon.

[![Build Status](https://travis-ci.org/alexellis/jaas.svg?branch=master)](https://travis-ci.org/alexellis/jaas)

**Motivation and context**

For a blog post covering use-cases for JaaS and more on the portions of the Docker API used see below:

* [Blog: One-shot containers on Docker Swarm](http://blog.alexellis.io/containers-on-swarm/)

**Related - Serverless**

If you would like to build Serverless applications with Docker Swarm or Kubernetes checkout my write-up on OpenFaaS:

* [Introducing Functions as a Service (FaaS)](https://blog.alexellis.io/introducing-functions-as-a-service/)

## Contributions are welcome

This is the contribution process for this repo.

* Raise a Github issue with the proposed change/idea
* I'll mark the issue as a feature/bug fix etc for the changelog
* This gives us a chance to discuss the idea
* If everything sounds good then go ahead and work on the PR
 * Please link to the bug and explain how you tested the change
* I'll merge after reviewing/testing

## Get started

### Build and install the code

Pre-requisites:

* Docker 1.13 or newer (experimental mode must be enabled if accessing service logs)
* [Go 1.7.3 (or Golang container)](https://golang.org/dl/)
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

* Run your first one-shot container:

```
# jaas -rm -image alexellis2/cows:latest
```

The `-rm` flag removes the Swarm service that was used to run your container. 

> The exit code from your container will also be available, you can check it with `echo $?`

* Hiding logs

If you aren't interested in the output logs then run it with the `--showlogs=false` override:

```
# jaas -image alexellis2/cows:latest --showlogs=false
```

* Removing service after completion

To remove the service after it completes, run with the `-rm` flag:

```
# jaas -image alexellis2/href-counter:latest --env url=http://blog.alexellis.io/ --showlogs=true

Service created: peaceful_shirley (uva6bcqyubm1b4c80dghjhb44)
ID:  uva6bcqyubm1b4c80dghjhb44  Update at:  2017-03-14 22:19:54.381973142 +0000 UTC
...

Exit code: 0
State: complete


Printing service logs
?2017-03-14T22:19:55.660902727Z com.docker.swarm.node.id=b2dqydhfavwezorhkqi11f962,com.docker.swarm.service.id=uva6bcqyubm1b4c80dghjhb44,com.docker.swarm.task.id=yruxuawdipz2v5n0wvvm8ib0r {"internal":42,"external":2}

Removing service...
```

* Using registryAuth

To enable pulling from secured registries you can use the `-registryAuth` parameter:
```
# export auth='{ "username" : "myUserName", "password" : "secret", "email" : "my@email", "serveraddress" : "my.reg.domain" }'
# export encAuth=`echo $auth | base64`
# jaas -registryAuth="$encAuth" -image my.reg.domain/hello-world:latest
```

*Notes on images*

You can have a multi-node swarm but make sure whatever image you choose is available in an accessible registry.

> A local image will not need to be pushed to a registry.

* Running jaas in a container

You can also run `jaas` in a container, but the syntax becomes slightly more verbose:

```
# docker run -ti -v /var/run/docker.sock:/var/run/docker.sock \
  alexellis2/jaas -image alexellis2/cows:latest
```

### Roadmap:

Here are several features / enhancements on the roadmap, please make additional suggestions through Github issues.

* [x] Optionally delete service after fetching exit code/logs
* [x] Support passing environmental variables
* [x] Support private registry auth via `-registryAuth` flag

Todo:

* [ ] Support constraints on where to run tasks
* [ ] Support optional secrets through CLI flag
* [ ] Validation around images which are not in local library
* [ ] Extract stdout/stderr etc from logs in human readable format similar to `docker logs`

### Future:

* When task logs are available in the API this will be used instead of service logs.
* When event streams are released they will prevent the need to poll continually
