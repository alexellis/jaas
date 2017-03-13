Readme for Raspberry Pi
=======================

This tool will work on the Raspberry Pi / an ARM v6/7 - either running with Go directly, or through Docker.

**Through Docker:**

Running:

```
$ docker run -ti -v /var/run/docker.sock:/var/run/docker.sock alexellis2/jaas-armhf ./jaas --showlogs=false -env url=http://www.alexellis.io -image=alexellis2/href-counter-armhf
```

Building:

```
$ docker build -t alexellis2/jaas-armhf . -f Dockerfile.armhf
```
