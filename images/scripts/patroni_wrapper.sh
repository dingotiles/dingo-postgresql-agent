#!/bin/bash

set -e
if [[ "${DEBUG}X" != "X" ]]; then
  set -x
fi

indent() {
  c="s/^/patroni> /"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

(
  # TODO: create /data/wal-e/env

  export PG_DATA_DIR=${DATA_VOLUME}/postgres0
  chown postgres:postgres -R ${DATA_VOLUME} /patroni /config

  if [[ -d ${PG_DATA_DIR} ]]; then
    chown postgres:postgres -R ${PG_DATA_DIR}
    chmod 700 $PG_DATA_DIR
  fi

  exec sudo PATH="${PATH}" -E -u postgres python3 /patroni/patroni.py /config/patroni.yml
) 2>&1 | indent
