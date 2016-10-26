#!/bin/bash

set -e -u

version=$(cat version/number)
digest=$(cat image/digest)

cat > dockerfile/Dockerfile << EOF
FROM dingotiles/dingo-postgresql@${digest}
ENV DINGO_IMAGE_VERSION=${version}
EOF
