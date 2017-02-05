#!/bin/bash

# TODO: not yet sure what function this hook script might perform for local-volume recovery of new leader

set -e # fail fast

indent() {
  c="s/^/restore_leader_if_missing> /"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

(
  echo "noop for local-volume backups"
) 2>&1 | indent
