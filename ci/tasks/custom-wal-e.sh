#!/bin/bash

set -e

cp -r agent/* dockerfile/

# remove submodule from development
rm -rf dockerfile/wal-e

# replace with wal-e to be tested
cp -r wal-e dockerfile/
