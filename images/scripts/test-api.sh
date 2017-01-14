#!/bin/sh

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
$DIR/initialize_etcd_auth.sh

dingo-postgresql-agent test-api
