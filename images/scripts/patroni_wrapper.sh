#!/bin/bash

set -e -u

indent() {
  c="s/^/patroni> /"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
patroni_config=/config/patroni.yml
patroni_env=/etc/patroni.d/.envrc

function wait_for_config {
  # wait for /config/patroni.yml to ensure that all variables stored in /etc/wal-e.d/env files
  wait_message="WARN: Waiting until ${patroni_env} and ${patroni_config} are created..."
  while [[ ! -f ${patroni_env} ]]; do
    if [[ "${wait_message}X" != "X" ]]; then
      echo ${wait_message} >&2
    fi
    sleep 1
    wait_message="" # only show wait_message once
  done
}

(
  if [[ "${DEBUG:-}X" != "X" ]]; then
    set -x


    echo "\nInstalled alpine/apk packages:"
    apk -vv info | sort

    echo "\nInstalled python3/pip3 packages":
    # ${DIR}/pip-versions.sh
    echo "TODO: re-enable pip-versions.sh"

  fi

  set +e
  wait_message="WARN: Waiting until ${patroni_config} is created..."
  if [[ ! -f ${patroni_config} ]]; then
    if [[ "${wait_message}X" != "X" ]]; then
      echo ${wait_message} >&2
    fi
    sleep 1
    wait_message=""
  fi
  set -e

  if [[ -d /etc/wal-e.d/env ]]; then
    echo "Environment variables provided to wal-e:"
    ls /etc/wal-e.d/env
  fi

  echo "/config/patroni.yml:"
  cat ${patroni_config}

  source ${patroni_env}
  env | sort

  : ${PATRONI_SCOPE:?required}
  : ${PG_DATA_DIR:?required}
  : ${ETCD_URI:?required}

  $DIR/restore_leader_if_missing.sh

  # runs as postgres user via supervisor
  python3 /patroni/patroni.py /config/patroni.yml
) 2>&1 | indent
