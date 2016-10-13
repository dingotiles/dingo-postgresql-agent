#!/bin/bash

package=$1
if [[ "${package}X" == "X" ]]; then
  function show_upload_time {
    package=$1
    version=$2
    upload_time=$(curl -s https://pypi.python.org/pypi/${package}/json | jq -r ".releases[\"${version}\"][-1].upload_time")
    echo $upload_time $package $version
  }

  pip3 list | sed -e "s/[()]//g" |
    while IFS= read -r package_info; do
      package=$(echo $package_info | awk '{print $1}')
      version=$(echo $package_info | awk '{print $2}')
      show_upload_time $package $version
    done
  exit
fi
