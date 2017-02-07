#!/bin/bash

set -e

: ${REMOTE_USER:?required}
: ${REMOTE_PUBLIC_KEY:?required}

if [[ ! -d /home/$REMOTE_USER ]]; then
  useradd -m -s /bin/bash $REMOTE_USER
fi

sshdir=/home/$REMOTE_USER/.ssh
mkdir -p ${sshdir}
echo -e $REMOTE_PUBLIC_KEY > ${sshdir}/authorized_keys

chown -R $REMOTE_USER:$REMOTE_USER ${sshdir}
chmod 700 ${sshdir}
chmod 600 ${sshdir}/authorized_keys

chown -R $REMOTE_USER:$REMOTE_USER /data

/usr/sbin/sshd -D
