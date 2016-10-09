# dingo-postgresql-agent

This WIP repo is to create a new agent to manage the configuration/lifecycle of Patroni; which in turn manages the configuration/lifecycle of PostgreSQL server.

To build:

```
docker build -t dingotiles/dingo-postgresql96-agent-base:latest .
```

Sanity check:

```
image=dingotiles/dingo-postgresql96-agent-base:latest
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
$ docker run -ti $image python3 --version
Python 3.5.2
```
