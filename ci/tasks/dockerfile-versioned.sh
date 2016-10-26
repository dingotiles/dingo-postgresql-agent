#!/bin/bash

set -e -u

version=$(cat version/number)
image_id=$(cat image/image-id)

cat > dockerfile/Dockerfile << EOF
FROM dingotiles/dingo-postgresql:${image_id}
ENV DINGO_IMAGE_VERSION=${VERSION}
EOF
