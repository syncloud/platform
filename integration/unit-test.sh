#!/usr/bin/env bash

APP_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )"  && cd .. && pwd )
cd ${APP_DIR}/src

echo "insalling test dependencies"
/opt/app/platform/python/bin/python /opt/app/platform/python/bin/pip2 install -U pytest
/opt/app/platform/python/bin/python /opt/app/platform/python/bin/pip2 install -r dev_requirements.txt

echo "installing in develop mode to run unit tests"
#pip2 install -e .
/opt/app/platform/python/bin/py.test --cov syncloud test
#python setup.py develop --uninstall