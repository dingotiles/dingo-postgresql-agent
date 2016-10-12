#!/usr/bin/dumb-init /bin/bash

# echo "Running any available init scripts:"
# for init in /config/init/*; do
#   echo " -- Running $init..."
#   $init
# done
#
echo "Starting supervisor"
supervisord
