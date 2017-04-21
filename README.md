# Ad-hoc Jobs as a Service (JaaS)

This project provides a simple Golang CLI tool that binds to the Docker Swarm API to create an ad-hoc/one-shot Service and then poll until it exits. Service logs can also be retrieved if the experimental feature is enabled on the Docker daemon.

[![Build Status](https://travis-ci.org/alexellis/jaas.svg?branch=master)](https://travis-ci.org/alexellis/jaas)

**Motivation and context**

For a blog post covering use-cases for JaaS and more on the portions of the Docker API used see below:

* [Blog: One-shot containers on Docker Swarm](http://blog.alexellis.io/containers-on-swarm/)

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

**Run these commands**

```
# export GOPATH=$HOME/go
# go get -d -v github.com/alexellis/jaas
# cd $GOPATH/src/github.com/alexellis/jaas
# go install
# export PATH=$PATH:$GOPATH/bin
```

Now test `jaas` with `jaas --help`

* Enable Swarm Mode

```
# docker swarm init
```

*Notes on images*

You can have a multi-node swarm but make sure whatever image you choose is available in an accessible registry.

> A local image will not need to be pushed to a registry.

### Running a task / batch job / one-shot container

* Run your first one-shot container:

```
# docker pull alexellis2/cows:latest
# jaas -rm -image alexellis2/cows:latest
```

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

* Running jaas in a container

You can also run `jaas` in a container, but the syntax becomes slightly more verbose:

```
# docker build -t jaas .
# docker run -ti -v /var/run/docker.sock:/var/run/docker.sock jaas -image alexellis2/cows:latest
```

### Roadmap:

Here are several features / enhancements on the roadmap, please make additional suggestions through Github issues.

* [x] Optionally delete service after fetching exit code/logs
* [x] Support passing environmental variables
* [ ] Extract stdout/stderr etc from logs in human readable format similar to `docker logs`
* [ ] Support optional secrets through CLI flag
* [ ] Validation around images which are not in local library

### Future:

* When task logs are available in the API this will be used instead of service logs.
* When event streams are released they will prevent the need to poll continually
