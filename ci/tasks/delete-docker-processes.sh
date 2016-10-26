#!/bin/sh

set -e -u

docker ps -a

docker ps -a | awk '{print $1}' | xargs -L1 docker rm -f
