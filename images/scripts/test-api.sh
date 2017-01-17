#!/bin/sh

set -e -x

echo Running test-api
env | sort

/scripts/initialize_etcd_auth.sh
dingo-postgresql-agent test-api
