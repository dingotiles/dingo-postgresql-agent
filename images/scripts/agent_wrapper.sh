#!/bin/bash

set -e
if [[ "${DEBUG}X" != "X" ]]; then
  set -x
fi

indent() {
  c="s/^/agent> /"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

(
  dingo-postgresql-agent run
) 2>&1 | indent
