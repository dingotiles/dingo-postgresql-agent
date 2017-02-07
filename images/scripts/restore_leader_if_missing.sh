#!/bin/bash

# restore_leader_if_missing.sh is a workaround for patroni not currently
# having a way to replicate a new leader from a wal-e backup.
#
# The idea is that patroni can create a replica from wal-e if there is a leader
# though I'm not sure why a leader is required. That's future work in patroni.
#
# So this script will create a fake leader to trick patroni into restoring a new
# leader from wal-e backup.
#
# It will only run the restoration process if:
# * there is no current leader,
# * if there is no local DB initialized, and
# * if there is a wal-e backup available

set -e # fail fast

indent() {
  c="s/^/restore_leader_if_missing> /"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

(
  if [[ "$(curl -s ${ETCD_URI:?required}/v2/keys/service/${PATRONI_SCOPE}/leader | jq -r .node.value)" != "null" ]]; then
    echo "leader exists, no additional preparation required for container to join cluster"
    exit 0
  fi
  if [[ -d ${PG_DATA_DIR}/global ]]; then
    echo "local database exists; no additional preparation required to restart container"
    exit 0
  fi

  echo Looking up existing backups:
  wal-e backup-list

  backups_lines=$(wal-e backup-list 2>/dev/null | wc -l)
  if [[ $backups_lines -lt 2 ]]; then
    echo "new cluster, no existing backup to restore"
    exit 0
  fi

  # must have /initialize set
  if [[ "$(curl -s ${ETCD_URI}/v2/keys/service/${PATRONI_SCOPE}/initialize | jq -r .node.value)" == "null" ]]; then
    if [[ "${WALE_S3_PREFIX:-X}" != "X" ]]; then
      echo "etcd missing /initialize system ID, fetching from ${WALE_S3_PREFIX:?required}sysids"
      region=$(aws s3api get-bucket-location --bucket ${WAL_S3_BUCKET} | jq -r '.LocationConstraint')
      if [[ ${region} != 'null' ]]; then
        region_option="--region ${region}"
      fi
      aws s3 ${region_option:-} sync ${WALE_S3_PREFIX}sysids /tmp/sysids
    elif [[ "${WALE_REMOTE_PREFIX:-X}" != "X" ]]; then
      mkdir -p /tmp/sysids
      echo "etcd missing /initialize system ID, fetching from ${WALE_REMOTE_PREFIX:?required}sysids/sysid"
      sysid=${REMOTE_BASE_PATH:?required}sysids/sysid
      set +e -x
      ssh -o StrictHostKeyChecking=no \
          -p ${REMOTE_PORT:-22} \
          -i ${REMOTE_IDENTITY_FILE:?required} \
          ${REMOTE_USER}@${REMOTE_HOST} \
          "[[ -f ${sysid} ]] && cat ${sysid}" > /tmp/sysids/sysid
      set -e +x
      echo "copied ${sysid} into /tmp/sysids/sysid: $(cat /tmp/sysids/sysid)"
    elif [[ "${WALE_LOCAL_PREFIX:-X}" != "X" ]]; then
      echo "etcd missing /initialize system ID, fetching from ${WALE_LOCAL_PREFIX:?required}sysids"
      cp -R ${LOCAL_BACKUP_VOLUME:?required}sysids /tmp/sysids
    else
      echo "Not implemented backup of sysids for '$ARCHIVE_METHOD'"
      exit 1
    fi


    if [[ ! -f /tmp/sysids/sysid ]]; then
      echo "Pre-existing backups but missing /sysids/sysid from original 'Database system identifier'"
      exit 1
    fi

    echo "Re-initializing /${PATRONI_SCOPE}/initialize with original 'Database system identifier'"
    curl -s ${ETCD_URI}/v2/keys/service/${PATRONI_SCOPE}/initialize -XPUT -d "value=$(cat /tmp/sysids/sysid)"
    curl -s ${ETCD_URI}/v2/keys/service/${PATRONI_SCOPE}/initialize
  fi

  echo "preparing patroni to restore this container from wal-e backups"
) 2>&1 | indent
