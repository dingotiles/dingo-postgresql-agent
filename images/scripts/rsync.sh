#!/bin/bash

set -e -u

path_name_of_file_to_archive=$2 # %p
file_name=$3                    # %f

case $1 in
  push )
    rsync -a $path_name_of_file_to_archive ${RSYNC_URI:?required}/wal_archive/$file_name
    ;;
  fetch )
    rsync -a ${RSYNC_URI}/wal_archive/$file_name $path_name_of_file_to_archive
    ;;
esac
