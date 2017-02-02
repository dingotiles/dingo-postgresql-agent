#!/bin/bash

# TODO: not yet sure what function this hook script might perform for rsync-based recovery of new leader

set -e # fail fast

indent() {
  c="s/^/restore_leader_if_missing> /"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

(
  echo "noop for rsync-based backups"
) 2>&1 | indent
