#!/bin/bash

set -u

indent() {
  c="s/^/base-backups> /"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

# $BACKUP_HOUR can be an hour in the day, or * to run backup each hour
BACKUP_HOUR=${BACKUP_HOUR:-1}
BACKUP_INTERVAL=${BACKUP_INTERVAL:-3600}

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PG_DATA_DIR=${DATA_VOLUME}/postgres0
patroni_env=/etc/patroni.d/.envrc

function wait_for_config {
  # wait for /config/patroni.yml to ensure that all variables stored in /etc/wal-e.d/env files
  wait_message="WARN: Waiting until ${patroni_env} is created..."
  while [[ ! -f ${patroni_env} ]]; do
    if [[ "${wait_message}X" != "X" ]]; then
      echo ${wait_message} >&2
    fi
    sleep 1
    wait_message="" # only show wait_message once
  done

  source ${patroni_env}
}

function wale_base_backups {
  # NOTE: env vars printed also ensures they are set (set -u)
  echo PATRONI_SCOPE: ${PATRONI_SCOPE}
  echo WALE_LOCAL_PREFIX: ${WALE_LOCAL_PREFIX:-}
  echo WALE_S3_PREFIX: ${WALE_S3_PREFIX:-}
  echo WAL_S3_BUCKET: ${WAL_S3_BUCKET:-}

  if [[ "${WALE_S3_PREFIX:-X}" != "X" ]]; then
    if [[ "${DEBUG:-X}" != "X" ]]; then
      aws s3api get-bucket-location --bucket ${WAL_S3_BUCKET:?required}
    fi
    AWS_REGION=$(aws s3api get-bucket-location --bucket ${WAL_S3_BUCKET:?required} | jq -r '.LocationConstraint')
    echo AWS_REGION: ${AWS_REGION}
  fi

  INITIAL=1
  SYSID_UPLOADED=0
  RETRY=0
  LAST_BACKUP_TS=0
  while true
  do
    sleep 5

    CURRENT_TS=$(date +%s)
    CURRENT_HOUR=$(date +%H)
    pg_isready >/dev/null 2>&2 || continue
    IN_RECOVERY=$(psql -tqAc "select pg_is_in_recovery()")

    if [[ $IN_RECOVERY != "f" ]]
    then
      echo "Not uploading backup because I am currently in recovery"

      if [[ $(du -s ${PG_DATA_DIR}/pg_xlog | awk '{print $1}') -gt 1000000 ]]
      then
        echo "WARNING pg_xlog dir is getting large (> 1GB)."

        # last_wal_file=$(pg_controldata ${PG_DATA_DIR} | grep "Latest checkpoint's REDO WAL file" | awk '{print $NF}')
        # if [[ "${last_wal_file}X" != "X" ]]; then
        #   find ${PG_DATA_DIR}/pg_xlog -maxdepth 1 -type f \
        #     \! -name $last_wal_file \
        #     \! -newer ${PG_DATA_DIR}/pg_xlog/$last_wal_file \
        #     -exec rm {} \;
        # fi
      fi

      continue
    fi

    # during initial run, count the number of backup lines. If there are
    # no backup (only line with backup-list header is returned), or there
    # is an error, try to produce a backup. Otherwise, stick to the regular
    # schedule, since we might run the backups on a freshly promoted replica.
    if [[ $INITIAL = 1 ]]
    then
      BACKUPS_LINES=$(wal-e backup-list 2>/dev/null|wc -l)
      [[ $PIPESTATUS[0] = 0 ]] && [[ $BACKUPS_LINES -ge 2 ]] && INITIAL=0
    fi
    # Backup the system ID
    if [[ $SYSID_UPLOADED = 0 ]]
    then
      pg_controldata ${PG_DATA_DIR}

      mkdir -p /tmp/sysids
      pg_controldata ${PG_DATA_DIR} | grep "Database system identifier" | cut -d ":" -f2 | awk '{print $1}' > /tmp/sysids/sysid

      if [[ "${WALE_S3_PREFIX:-X}" != "X" ]]; then
        if [[ ${AWS_REGION} != 'null' ]]; then
          region_option="--region ${AWS_REGION}"
        fi
        aws s3 ${region_option:-} sync /tmp/sysids ${WALE_S3_PREFIX}sysids
      elif [[ "${WALE_SSH_PREFIX:-X}" != "X" ]]; then
        sysids=${SSH_BASE_PATH:?required}sysids
        ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
            -p ${SSH_PORT:-22} \
            -i ${SSH_IDENTITY_FILE:?required} \
            ${SSH_USER}@${SSH_HOST} \
            "mkdir -p ${sysids} && cat > ${sysids}/sysid" < /tmp/sysids/sysid
      elif [[ "${WALE_LOCAL_PREFIX:-X}" != "X" ]]; then
        cp -R /tmp/sysids ${LOCAL_BACKUP_VOLUME:?required}sysids
      else
        echo "Not implemented backup of sysids for '$ARCHIVE_METHOD'"
        exit 1
      fi
      SYSID_UPLOADED=1
    fi
    # produce backup only at a given hour, unless it's set to *, which means
    # that only backup_interval is taken into account. We also skip all checks
    # when the backup is forced because of previous attempt's failure or because
    # it's going to be a very first backup, in which case we create it unconditionally.
    if [[ $RETRY = 0 ]] && [[ $INITIAL = 0 ]]
    then
      # check that enough time has passed since the previous backup
      [[ $BACKUP_HOUR != '*' ]] && [[ $CURRENT_HOUR != $BACKUP_HOUR ]] && continue
      # get the time since the last backup. Do it only one when the hour
      # matches the backup hour.
      [[ $LAST_BACKUP_TS = 0 ]] && LAST_BACKUP_TS=$(wal-e backup-list LATEST 2>/dev/null | tail -n1 | awk '{print $2}' | xargs date +%s --date)
      # LAST_BACKUP_TS will be empty on error.
      if [[ -z $LAST_BACKUP_TS ]]
      then
        LAST_BACKUP_TS=0
        echo "could not obtain latest backup timestamp"
      fi

      ELAPSED_TIME=$((CURRENT_TS-LAST_BACKUP_TS))
      [[ $ELAPSED_TIME -lt $BACKUP_INTERVAL ]] && continue
    fi
    # leave only 4 base backups before creating a new one
    wal-e delete --confirm retain 4
    # push a new base backup
    echo "producing a new backup at $(date)"
    # reduce the priority of the backup for CPU consumption
    PGUSER=${REPLICATION_USER} nice -n 5 wal-e backup-push ${PG_DATA_DIR}
    RETRY=$?
    # re-examine last backup timestamp if a new backup has been created
    if [[ $RETRY = 0 ]]
    then
      INITIAL=0
      LAST_BACKUP_TS=0
    fi
  done
}

(
  echo Waiting for configuration from agent...
  wait_for_config
  echo Configuration acquired from agent, beginning loop for base backups...

  source $patroni_env
  if [[ "${DEBUG:-}X" != "X" ]]; then
    env | sort
  fi

  wale_base_backups
) 2>&1 | indent
