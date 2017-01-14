#!/bin/sh

set -e

/scripts/initialize_etcd_auth.sh

dingo-postgresql-agent test-api
