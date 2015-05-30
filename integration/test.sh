#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

echo "rootfs version (docker):$(</version)"
pip2 freeze | grep syncloud
pip2 install -U pytest
pip2 install -r /requirements.txt
python setup.py
syncloud-platform-post-install
export TEAMCITY_VERSION=9
py.test -s verify.py