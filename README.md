# dingo-postgresql-agent

**PRIVATE: active AWS credentials in `.envrc` that need to be rotated before making repo public**
**PRIVATE: still in pre-alpha development; will be OSS later**

This WIP repo is to create a new agent to manage the configuration/lifecycle of Patroni; which in turn manages the configuration/lifecycle of PostgreSQL server.

## Proposed demo for homepage

```
# local docker
DOCKER_HOST_IP=$(ifconfig | grep inet | grep -v inet6 | grep -v 127.0.0.1 | head -n1 | awk '{print $2}')
# docker-machine
DOCKER_HOST_IP=$(docker-machine ip <name>)

docker run -ti \
  --name dingo-demo \
  -e DINGO_CLUSTER=dingo-demo \
  -e DINGO_ORG_TOKEN=shared-org \
  -e DOCKER_HOST_IP=${DOCKER_HOST_IP} \
  -e DOCKER_HOST_PORT_5432=5000 \
  -v 5000:5432 \
  dingotiles/dingo-postgresql96-agent /scripts/entry.sh
```

You can poll to check when PostgreSQL is running and is successfully continuously archiving itself:

```
uri=postgres://superuser:password@${DOCKER_HOST_IP}:5000/postgres
psql $url -c "SELECT current_database();"
```

**TODO:** Your demonstration PostgreSQL has default passwords. That's ok because you can change passwords - an important feature to allow you to rotate credentials throughout the life of the database.




## Sample cluster

The project includes a sample `docker-compose.yml` to run a sample cluster of an agent, test-api, registrator and etcd:

```
DOCKER_HOST_IP=$(docker-machine ip <machine-name>)
docker-compose up
```

## Development

To build:

```
docker build -t dingotiles/dingo-postgresql96-agent-base:latest images/pg96-base
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
