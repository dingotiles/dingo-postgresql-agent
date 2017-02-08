#!/bin/bash

set -e

if [[ "${DEBUG:-}X" != "X" ]]; then
  set -x
fi

indent() {
  c="s/^/agent> /"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

(
mkdir -p ${DATA_VOLUME} /patroni /config /etc/patroni.d/
chown postgres:postgres -R ${DATA_VOLUME} /patroni /config /etc/patroni.d/

export PG_DATA_DIR=${DATA_VOLUME}/postgres0

if [[ -d ${PG_DATA_DIR} ]]; then
  chown postgres:postgres -R ${PG_DATA_DIR}
  chmod 700 $PG_DATA_DIR
fi

if [[ "${LOCAL_BACKUP_VOLUME:-X}" != "X" ]]; then
  mkdir -p $LOCAL_BACKUP_VOLUME
  chown -R postgres:postgres $LOCAL_BACKUP_VOLUME
  chmod 700 $LOCAL_BACKUP_VOLUME
fi

mkdir -p /home/postgres/.ssh
touch /home/postgres/.ssh/ssh_backup_storage
chown -R postgres:postgres /home/postgres
chmod 700 /home/postgres/.ssh
chmod 600 /home/postgres/.ssh/ssh_backup_storage
# file contents come from API

export PATRONI_POSTGRES_START_COMMAND="supervisorctl start postgres:"

dingo-postgresql-agent run
) 2>&1 | indent
