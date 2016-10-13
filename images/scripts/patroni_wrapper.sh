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

  if [[ "${DEBUG}X" != "X" ]]; then
    set -x

    DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

    echo "\nInstalled alpine/apk packages:"
    apk info | xargs -I % apk info % | grep description: | awk '{print $1}' | sort

    echo "\nInstalled python3/pip3 packages":
    ${DIR}/pip-versions.sh

  fi

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

  export PG_DATA_DIR=${DATA_VOLUME}/postgres0
  chown postgres:postgres -R ${DATA_VOLUME} /patroni /config

  if [[ -d ${PG_DATA_DIR} ]]; then
    chown postgres:postgres -R ${PG_DATA_DIR}
    chmod 700 $PG_DATA_DIR
  fi

  exec sudo PATH="${PATH}" -E -u postgres python3 /patroni/patroni.py /config/patroni.yml
) 2>&1 | indent
