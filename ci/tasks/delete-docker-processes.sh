#!/bin/sh

set -e -u

env
ls -al /root/.docker/machine/machines/
ls -al /root/.docker/machine/machines/$MACHINE_NAME

docker-machine ls
exit 1

DOCKER_VERSION=1.12.2
curl -L https://get.docker.com/builds/Linux/x86_64/docker-${DOCKER_VERSION}.tgz -o docker-${DOCKER_VERSION}.tgz
tar xfz docker-${DOCKER_VERSION}.tgz
mv docker/docker /usr/local/bin/
rm -rf docker

eval $(docker-machine env $MACHINE_NAME)
env

docker ps -a

docker ps -a | awk '{print $1}' | xargs -L1 docker rm -f
