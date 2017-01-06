# jaas
Jobs as a Service (JaaS) for Docker Swarm

## Get started

* Pre-requisites:
* Docker 1.13-RC
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

You can also run alexellis2/jaas in a container, but the syntax becomes slightly more verbose:

```
# docker build -t jaas .
# docker run -ti -v /var/run/docker.sock:/var/run/docker.sock jaas -image alexellis2/cows:latest
```
