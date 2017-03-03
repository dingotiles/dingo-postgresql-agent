#!/bin/bash

set -u

indent() {
  c="s/^/backups-summary> /"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

STARTUP_POLL_WAITTIME=${STARTUP_POLL_WAITTIME:-5}
BACKUPS_SUMMARY_WAITTIME=${BACKUPS_SUMMARY_WAITTIME:-60}

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
patroni_env=/etc/patroni.d/.envrc

function wait_for_config {
  # TODO: wait for "supervisorctl status agent" to be RUNNING
  # wait for /config/patroni.yml to ensure that all variables stored in /etc/wal-e.d/env files
  wait_message="WARN: Waiting until ${patroni_env} is created..."
  while [[ ! -f ${patroni_env} ]]; do
    if [[ "${wait_message}X" != "X" ]]; then
      echo ${wait_message} >&2
    fi
    sleep 1
    wait_message="" # only show wait_message once
  done
}

source ${patroni_env}
echo ETCD_CLUSTER_URI=${ETCD_CLUSTER_URI:?required}

backups_found=
function backups_summary {
  curl -s ${ETCD_CLUSTER_URI}/wale-backup-list \
    -X PUT -d "value=$(wal-e backup-list 2>/dev/null)" > /dev/null
  backup_lines=$(wal-e backup-list 2>/dev/null | wc -l)
  if [[ $backup_lines -ge 2 ]]; then
    if [[ "${DEBUG:-}X" != "X" ]]; then
      echo "INFO: Backup status:"
      wal-e backup-list 2>/dev/null
    fi
    backups_found=1
  else
    echo "WARNING: No backups successful yet"
    backups_found=
  fi
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
    pg_isready >/dev/null 2>&2 || {
      echo "database is not yet ready for connections: pg_isready"
      sleep ${STARTUP_POLL_WAITTIME}
      continue
    }
    in_recovery=$(psql -U postgres -tqAc "select pg_is_in_recovery()")
    if [[ "${in_recovery}" != "f" ]]; then
      echo "database is still in recovery: select pg_is_in_recovery()"
      sleep ${STARTUP_POLL_WAITTIME}
      continue
    fi
    backups_summary
    if [[ "${backups_found}" == "1" ]]; then
      sleep ${BACKUPS_SUMMARY_WAITTIME}
    else
      sleep ${STARTUP_POLL_WAITTIME}
    fi
  done
) 2>&1 | indent
