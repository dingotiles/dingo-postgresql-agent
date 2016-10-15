#!/bin/bash

indent() {
  c="s/^/backups-summary> /"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

wale_env_dir=/etc/wal-e.d/env
patroni_config=/config/patroni.yml
BACKUPS_SUMMARY_WAITTIME=${BACKUPS_SUMMARY_WAITTIME:-60}

function wait_for_config {
  # wait for /config/patroni.yml to ensure that all variables stored in /etc/wal-e.d/env files
  wait_message="WARN: Waiting until ${wale_env_dir} and ${patroni_config} are created..."
  if [[ ! -d ${wale_env_dir} || ! -f ${patroni_config} ]]; then
    if [[ "${wait_message}X" != "X" ]]; then
      echo ${wait_message} >&2
    fi
    sleep 1
    wait_message="" # only show wait_message once
  fi
}

function backups-list {
  wal-e backup-list
  # wal-e backup-list 2>/dev/null
  curl -s ${ETCD_HOST_PORT}/v2/keys/service/${PATRONI_SCOPE}/wale-backup-list \
    -X PUT -d "value=$(wal-e backup-list 2>/dev/null)" > /dev/null
}

(
  echo Waiting for configuration from agent...
  wait_for_config
  echo Configuration acquired from agent, beginning loop to dump backup summaries...

  env | sort

  while true; do
    in_recovery=$(psql -tqAc "select pg_is_in_recovery()")
    if [[ "${in_recovery}" == "f" ]]; then
      envdir wale_env_dir backups-list
    fi
    sleep ${BACKUPS_SUMMARY_WAITTIME}
  done
) 2>&1 | indent
