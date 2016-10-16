#!/bin/bash

set -e

indent() {
  c="s/^/patroni> /"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

(
  wale_env_dir=/etc/wal-e.d/env
  patroni_config=/config/patroni.yml

  if [[ "${DEBUG:-}X" != "X" ]]; then
    set -x

    DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

    echo "\nInstalled alpine/apk packages:"
    apk -vv info | sort

    echo "\nInstalled python3/pip3 packages":
    ${DIR}/pip-versions.sh

  fi

  set +e
  wait_message="WARN: Waiting until ${wale_env_dir} is created..."
  if [[ ! -d ${wale_env_dir} ]]; then
    if [[ "${wait_message}X" != "X" ]]; then
      echo ${wait_message} >&2
    fi
    sleep 1
    wait_message=""
  fi

  echo "Environment variables provided to wal-e:"
  ls ${wale_env_dir}

  wait_message="WARN: Waiting until ${patroni_config} is created..."
  if [[ ! -f ${patroni_config} ]]; then
    if [[ "${wait_message}X" != "X" ]]; then
      echo ${wait_message} >&2
    fi
    sleep 1
    wait_message=""
  fi
  set -e

  echo "/config/patroni.yml:"
  cat ${patroni_config}

  env | sort

  # runs as postgres user via supervisor
  python3 /patroni/patroni.py /config/patroni.yml
) 2>&1 | indent
