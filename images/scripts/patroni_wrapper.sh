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

  if [[ -f /etc/wal-e.d/env/SSH_PRIVATE_KEY ]]; then
    echo -e $(cat /etc/wal-e.d/env/SSH_PRIVATE_KEY) > /home/postgres/.ssh/ssh_backup_storage
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
