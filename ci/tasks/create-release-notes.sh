#!/bin/bash

set -e

NOTES=$PWD/release-notes

version=$(cat version/number)

echo v${version} > $NOTES/release-name

cat > $NOTES/notes.md <<EOF
## Improvements

TODO

## Installation

Normal image (with [Postgis](http://postgis.net/) extension support):

\`\`\`
docker pull dingotiles/dingo-postgresql:latest
docker pull dingotiles/dingo-postgresql:v${version}

docker run dingotiles/dingo-postgresql
\`\`\`

Image with [Madlib](https://madlib.incubator.apache.org/) extension support:

\`\`\`
docker pull dingotiles/dingo-postgresql-madlib:latest
docker pull dingotiles/dingo-postgresql-madlib:v${version}

docker run dingotiles/dingo-postgresql-madlib
\`\`\`
EOF
