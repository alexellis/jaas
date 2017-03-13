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

### Pre-requisites:

* Docker 1.13 or newer (experimental mode must be enabled if accessing service logs)
* Go 1.7.3 (or Golang container)

* Build the code:

```
# go install
```

* Enable Swarm Mode

```
# docker swarm init
```

You can have a multi-node swarm but make sure whatever image you choose is available in an accessible registry.

> A local image will not need to be pushed to a registry.

* Run your first one-shot container:

```
# docker pull alexellis2/cows:latest
# jaas -image alexellis2/cows:latest
```

* Hiding logs

If you aren't interested in the output logs then run it with the `--showlogs=false` override:

```
# jaas -image alexellis2/cows:latest --showlogs=false
```

* Running jaas in a container

You can also run `jaas` in a container, but the syntax becomes slightly more verbose:

```
# docker build -t jaas .
# docker run -ti -v /var/run/docker.sock:/var/run/docker.sock jaas -image alexellis2/cows:latest
```

### Roadmap:

Here are several features / enhancements on the roadmap, please make additional suggestions through Github issues.

* [ ] Extract stdout/stderr etc from logs in human readable format similar to `docker logs`
* [ ] Optionally delete service after fetching exit code/logs
* [ ] Support optional secrets through CLI flag
* [ ] Validation around images which are not in local library
* [x] Support passing environmental variables

### Future:

* When task logs are available in the API this will be used instead of service logs.
* When event streams are released they will prevent the need to poll continually
