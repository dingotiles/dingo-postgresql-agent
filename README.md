# dingo-postgresql-agent

This WIP repo is to create a new agent to manage the configuration/lifecycle of Patroni; which in turn manages the configuration/lifecycle of PostgreSQL server.

To build:

```
docker build -t dingotiles/dingo-postgresql96-agent:latest .
```

Sanity check:

```
image=dingotiles/dingo-postgresql96-agent:latest
$ docker run -ti $image dingo-postgresql-agent --version
dingo-postgresql-agent version 0.1.0

$ docker run -ti $image postgres --version
postgres (PostgreSQL) 9.6.0
```

See images/pg96-base/README.md for additional sanity checking of contents.
