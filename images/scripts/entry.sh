#!/usr/bin/dumb-init /bin/bash

set -e

/scripts/show_usage.sh

echo "Starting supervisor"
supervisord -c /etc/supervisor/supervisord.conf
