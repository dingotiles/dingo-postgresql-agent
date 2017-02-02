#!/bin/bash

indent() {
  c="s/^/create_replica> /"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

(

  echo "Running create_replica_method for rsync"

  env

  echo args $@
) 2>&1 | indent
