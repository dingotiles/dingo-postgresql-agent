#!/bin/bash

# Built from https://github.com/nabeken/docker-volume-container-rsync
# USAGE:
#   rsync rsync://${DOCKER_HOST_IP}:4873/
#   rsync rsync://${DOCKER_HOST_IP}:4873/volume/
# To upload/sync files:
#   rsync -aP images/backup_storage rsync://${DOCKER_HOST_IP}:4873/volume/

VOLUME=${VOLUME:-/data}
OWNER=${OWNER:-nobody}
GROUP=${GROUP:-nogroup}

chown "${OWNER}:${GROUP}" "${VOLUME}"

[ -f /etc/rsyncd.conf ] || cat <<EOF > /etc/rsyncd.conf
uid = ${OWNER}
gid = ${GROUP}
use chroot = yes
log file = /dev/stdout
reverse lookup = no
[volume]
    hosts allow = *
    read only = false
    path = ${VOLUME}
    comment = docker volume
EOF

exec /usr/bin/rsync --no-detach --daemon --config /etc/rsyncd.conf "$@"
