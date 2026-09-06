#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

npm config ls -l
npm config set fetch-retry-mintimeout 200000
npm config set fetch-retry-maxtimeout 1200000

"${DIR}/platform/build.sh"
"${DIR}/login/build.sh"
