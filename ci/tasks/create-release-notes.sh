#!/bin/bash

set -e

NOTES=$PWD/release-notes

version=$(cat version/number)

echo v${version} > $NOTES/release-name

cat > $NOTES/notes.md <<EOF
## Improvements

TODO

## Upgrade

```
docker pull dingotiles/dingo-postgresql:latest
docker pull dingotiles/dingo-postgresql:v${version}
```
EOF
