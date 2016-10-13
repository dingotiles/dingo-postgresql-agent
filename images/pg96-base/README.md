# dingo-postgresql96-agent-base Docker image

To see the versioned aspects of the base Docker image, run:

```
grep "^ENV.*VERSION" images/pg96-base/Dockerfile
```

The output will look similar to:

```
ENV PG_VERSION=9.6
ENV POSTGIS_VERSION 2.2.2
ENV WALE_VERSION=1.0.0
ENV SUPERVISOR_VERSION=3.3.1
ENV GOSU_VERSION=1.10
ENV PATRONI_VERSION=1.1
```

To build:

```
docker build -t dingotiles/dingo-postgresql96-agent-base:latest images/pg96-base
```

Sanity check of installed tools:

```
$ image=dingotiles/dingo-postgresql96-agent-base:latest

$ docker run -ti $image postgres --version
postgres (PostgreSQL) 9.6.0
$ docker run -ti $image psql --version
psql (PostgreSQL) 9.5.4
$ docker run -ti $image jq --version
jq-1.5
$ docker run -ti $image etcdctl --version
etcdctl version 2.3.4
$ docker run -ti $image wal-e version
1.0.0
$ docker run -ti $image python --version
Python 2.7.12
$ docker run -ti $image python3 --version
Python 3.5.2
$ docker run -ti $image gosu | grep version
gosu version: 1.10 (go1.7.1 on linux/amd64; gc)
$ docker run -ti $image dumb-init --version
dumb-init v1.2.0
$ docker run -ti $image python3 /patroni/patroni.py
Usage: /patroni/patroni.py config.yml
       	Patroni may also read the configuration from the PATRONI_CONFIGURATION environment variable
```

NOTE: python (2.7) is only required for supervisor. Supervisor 4.0+ will support python3 and then we can remove python (2.7) from base image.
