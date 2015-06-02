#!/bin/bash

APP_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
cd ${APP_DIR}

export TEAMCITY_VERSION=9

#Fix debconf frontend warnings
export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
export DEBCONF_FRONTEND=noninteractive
export DEBIAN_FRONTEND=noninteractive

echo "=== install syncloud dependencies ==="
wget --no-check-certificate --progress=dot:mega -O get-pip.py https://bootstrap.pypa.io/get-pip.py 2>&1
python get-pip.py

pip2 install -U pytest
pip2 install -r dev_requirements.txt

pip2 install -e .
py.test --cov syncloud test
python setup.py develop --uninstall

echo "installing python part"
python setup.py sdist
pip2 install --no-binary :all: dist/syncloud-platform-*.tar.gz
pip2 freeze | grep syncloud


echo "installing binary part"
syncloud-platform-post-install platform.tar.gz

py.test -s integration/verify.py