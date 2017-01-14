#!/bin/bash

set -e

# This script is used by tests to setup authorization user on etcd

if [[ "$(curl -s ${ETCD_HOST_PORT:?required}/v2/auth/users |  jq -r .message)" != "Insufficient credentials" ]]; then
  if [[ "$(curl -s ${ETCD_HOST_PORT:?required}/v2/auth/users |  jq -r .users)" == "null" ]]; then
    echo "Creating missing root user..."
    curl -s ${ETCD_HOST_PORT:?required}/v2/auth/users/root -X PUT -d "{\"user\":\"${ETCD_USERNAME:-root}\",\"password\":\"${ETCD_PASSWORD:?required}\"}"
  fi

  auth_enabled=$(curl -s ${ETCD_HOST_PORT:?required}/v2/auth/enable | jq -r .enabled)
  if [[ "${auth_enabled}" != "true" ]]; then
    echo "Enabling etcd authorization..."
    curl -s ${ETCD_HOST_PORT:?required}/v2/auth/enable -X PUT
  fi
fi

echo "Verifying crdentials..."
curl -s -u ${ETCD_USERNAME:-root}:${ETCD_PASSWORD:?required} ${ETCD_HOST_PORT}/v2/auth/users
