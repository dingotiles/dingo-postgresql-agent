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

if [[ -d /etc/wal-e.d/env ]]; then
  envdir /etc/wal-e.d/env wal-e wal-fetch "$file_name" "$path_name_of_file_to_archive" -p 1
else
  (>&2 echo "restore_command.sh has not been provided /etc/wal-e.d/env, exiting...")
  exit 1
fi
