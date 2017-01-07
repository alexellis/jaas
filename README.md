# Ad-hoc Jobs as a Service (JaaS)

This project provides a simple CLI that binds to the Docker Swarm API to create an ad-hoc/one-shot Service and then poll until it exits. Service logs can also be retrieved if the experimental feature is enabled on your Docker Engine.

## Get started

### Pre-requisites:

* Docker 1.13-RC (in experimental mode for service logs)
* Go 1.7.3 (or Golang container)

* Build the code:

```
# go build
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
# ./jaas -image alexellis2/cows:latest
```

* Hiding logs

If you aren't interested in the output logs then run it with the `--showlogs=false` override:

```
# ./jaas -image alexellis2/cows:latest --showlogs=false
```

* Running jaas in a container

You can also run `jaas` in a container, but the syntax becomes slightly more verbose:

```
# docker build -t jaas .
# docker run -ti -v /var/run/docker.sock:/var/run/docker.sock jaas -image alexellis2/cows:latest
```

### Todo:

I'd like suggestions on what else you need to make this usable.

* Optionally delete service after fetching exit code/logs
* Support optional secrets through CLI flag
* Validation around images which are not in local library

### Future:

* When task logs are available in the API this will be used instead of service logs.
* When event streams are released they will prevent the need to poll continually
