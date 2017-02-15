#!/bin/bash

function show_usage {
  guid=$(slug)
  sample_cluster_name=demo-cluster-${guid:0:6}
  cat <<'USAGE'
 ____  _
|  _ \(_)_ __   __ _  ___
| | | | | '_ \ / _` |/ _ \
| |_| | | | | | (_| | (_) |
|____/|_|_| |_|\__, |\___/
              |___/

Dingo PostgreSQL continuously archived from Day 1.

How to run:

  docker run -d \
    --name dingo-postgresql \
    -e DOCKER_HOST_IP=${DOCKER_HOST_IP:?required} \
    -e DOCKER_HOST_PORT_5432=${PUBLIC_PORT:-5000} \
    -p ${PUBLIC_PORT:-5000}:5432 \
    -e DINGO_ACCOUNT=global-org \
USAGE
echo "    -e DINGO_CLUSTER=${sample_cluster_name} \\"
cat <<'USAGE'
    dingotiles/dingo-postgresql

How to get direct PostgreSQL URI:

  uri=$(docker exec dingo-postgresql cat /config/uri)
  psql $uri
USAGE
}

if [[ -z "${DOCKER_HOST_IP:+x}" || -z "${DOCKER_HOST_PORT_5432:+x}" || \
      -z "${DINGO_CLUSTER:+x}" || -z "${DINGO_ACCOUNT:+x}" ]]; then
  show_usage

  echo
  echo "Required environment variables:"
  echo "  DOCKER_HOST_IP=$DOCKER_HOST_IP"
  echo "  DOCKER_HOST_PORT_5432=$DOCKER_HOST_PORT_5432"
  echo "  DINGO_CLUSTER=$DINGO_CLUSTER"
  echo "  DINGO_ACCOUNT=$DINGO_ACCOUNT"
  exit 1
fi
