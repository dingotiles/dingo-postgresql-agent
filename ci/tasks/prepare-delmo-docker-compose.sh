#!/bin/bash

set -e # fail fast
set -x # print commands

: ${image_name_tag:?required}

git clone agent agent-delmo
cd agent-delmo
mkdir -p tmp

cat > tmp/docker-compose-overrides.yml <<YAML
---
services:
  patroni1:
    image: ${image_name_tag}
YAML

spruce merge docker-compose.yml tmp/docker-compose-overrides.yml > tmp/docker-compose.yml
cp tmp/docker-compose.yml docker-compose.yml
