#!/bin/bash

# USAGE in patroni config: restore_command: "/scripts/restore_command.sh \"%p\" \"%f\""
#
# From documentation:
# In archive_command:
# - %p is replaced by the path name of the file to archive,
# - while %f is replaced by only the file name.
# (The path name is relative to the current working directory, i.e., the cluster's data directory.)
# Use %% if you need to embed an actual % character in the command.

set -e -u

path_name_of_file_to_archive=$1 # %p
file_name=$2                    # %f

if [[ "${WALE_S3_ENDPOINT:-X}" != "X" ]]; then
  envdir /etc/wal-e.d/env wal-e wal-fetch \"%f\" \"${}\" -p 1
elif [[ "${RSYNC_HOSTNAME:-X}" != "X" ]]; then
  rsync -a ${RSYNC_USERNAME}@${RSYNC_HOSTNAME}:${RSYNC_DEST_DIR}/wal_archive/$file_name $path_name_of_file_to_archive
else
  (>&2 echo "archive_command.sh has not been provided \$WALE_S3_ENDPOINT nor \$RSYNC_HOSTNAME, exiting...")
  exit 1
fi
