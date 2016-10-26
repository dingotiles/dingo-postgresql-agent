#!/bin/sh

set -e -u

DOCKER_VERSION=1.12.2
mkdir /tmp/
curl -L https://get.docker.com/builds/Linux/x86_64/docker-${DOCKER_VERSION}.tgz -o /tmp/docker-${DOCKER_VERSION}.tgz
tar xfz /tmp/docker-${DOCKER_VERSION}.tgz /tmp/
mv /tmp/docker/docker /usr/local/bin/
rm -rf /tmp/docker

docker ps -a

docker ps -a | awk '{print $1}' | xargs -L1 docker rm -f
