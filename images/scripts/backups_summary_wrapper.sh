#!/bin/bash

set -u

indent() {
  c="s/^/backups-summary> /"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

BACKUPS_SUMMARY_WAITTIME=${BACKUPS_SUMMARY_WAITTIME:-60}

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PG_DATA_DIR=${DATA_VOLUME}/postgres0
patroni_env=/etc/patroni.d/.envrc

function wait_for_config {
  # wait for /config/patroni.yml to ensure that all variables stored in /etc/wal-e.d/env files
  wait_message="WARN: Waiting until ${patroni_env} are created..."
  if [[ ! -f ${patroni_env} ]]; then
    if [[ "${wait_message}X" != "X" ]]; then
      echo ${wait_message} >&2
    fi
    sleep 1
    wait_message="" # only show wait_message once
  fi
}

function backups_summary {
  # display variables & enforce (set -u) that they are already set
  echo PATRONI_SCOPE: ${PATRONI_SCOPE}
  echo ETCD_HOST_PORT: ${ETCD_HOST_PORT}
  echo WALE_S3_PREFIX: ${WALE_S3_PREFIX}
  echo WAL_S3_BUCKET: ${WAL_S3_BUCKET}

  wal-e backup-list 2>/dev/null
  curl -s ${ETCD_HOST_PORT}/v2/keys/service/${PATRONI_SCOPE}/wale-backup-list \
    -X PUT -d "value=$(wal-e backup-list 2>/dev/null)" > /dev/null
}

(
  echo Waiting for configuration from agent...
  wait_for_config
  echo Configuration acquired from agent, beginning loop for base backups...

  source $patroni_env

  if [[ "${DEBUG:-}X" != "X" ]]; then
    env | sort
  fi

  while true; do
    pg_isready >/dev/null 2>&2 || continue
    in_recovery=$(psql -tqAc "select pg_is_in_recovery()")
    echo "in_recovery: ${in_recovery}"
    if [[ "${in_recovery}" == "f" ]]; then
      backups_summary
    fi
    sleep ${BACKUPS_SUMMARY_WAITTIME}
  done
) 2>&1 | indent
