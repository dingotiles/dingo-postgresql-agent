# dingo-postgresql96-agent-base Docker image

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
