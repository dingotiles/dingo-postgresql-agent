#!/bin/bash

set -e

indent() {
  c="s/^/backups> /"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

(
  echo "TODO: base backups"
) 2>&1 | indent
