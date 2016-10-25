#!/bin/bash

set -e -u

VERSION=$(cat version/number)

cat > dockerfile/Dockerfile << EOF
FROM dingotiles/dingo-postgresql:latest
ENV DINGO_IMAGE_VERSION=${VERSION}
EOF
