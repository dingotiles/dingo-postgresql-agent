#!/bin/bash

indent() {
  c="s/^/advertise> /"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

patroni_env=/etc/patroni.d/.envrc

(
  set +e
  # patroni_env is created by the agent during initialization
  wait_message="Waiting until ${patroni_env} is created..."
  while [[ ! -f ${patroni_env} ]]; do
    if [[ "${wait_message}X" != "X" ]]; then
      echo ${wait_message} >&2
    fi
    sleep 1
    wait_message=""
  done
  set -e

  source ${patroni_env}

  echo "Self-advertising cell/node pairing for cluster"
  echo DINGO_NODE=${DINGO_NODE:?required}
  echo DINGO_CLUSTER=${DINGO_CLUSTER:?required}
  echo ETCD_CLUSTER_URI=${ETCD_CLUSTER_URI:?required}
  echo DOCKER_HOST_IP=${DOCKER_HOST_IP:?required}
  CELL_GUID=${CELL_GUID:-$DOCKER_HOST_IP}
  echo CELL_GUID=${CELL_GUID}

  while true; do
    value=$( \
      curl -s localhost:8008 | \
      jq -c \
        --arg cell ${CELL_GUID} \
        --arg node ${DINGO_NODE} \
        '{cell_guid:$cell, node_id:$node, state:.state, role:.role}' \
      )

    if [[ "${value}X" == "X" ]]; then
      value="{\"cell_guid\":\"${CELL_GUID}\",\"node_id\":\"${DINGO_NODE}\",\"state\":\"api-not-available\"}"
      echo value=$value
    fi
    curl -sf ${ETCD_CLUSTER_URI}/nodes/${DINGO_NODE}?ttl=20 \
      -XPUT -d "value=${value}" >/dev/null

    sleep 6
  done
) 2>&1 | indent
