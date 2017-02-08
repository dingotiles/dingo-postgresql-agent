#!/bin/bash

set -e

: ${SSH_USER:?required}
: ${SSH_PUBLIC_KEY:?required}

if [[ ! -d /home/$SSH_USER ]]; then
  useradd -m -s /bin/bash $SSH_USER
fi

sshdir=/home/$SSH_USER/.ssh
mkdir -p ${sshdir}
echo -e $SSH_PUBLIC_KEY > ${sshdir}/authorized_keys

chown -R $SSH_USER:$SSH_USER ${sshdir}
chmod 700 ${sshdir}
chmod 600 ${sshdir}/authorized_keys

chown -R $SSH_USER:$SSH_USER /data

/usr/sbin/sshd -D
