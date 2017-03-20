#!/bin/bash

set -e -u

version=$(cat version/number)
digest=$(cat image/digest)
: ${image_name:?required}

cat > dockerfile/Dockerfile << EOF
FROM ${image_name}@${digest}
ENV DINGO_IMAGE_VERSION=${version}
EOF
