#!/usr/bin/env bash

APP_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )"  && pwd )
cd ${APP_DIR}

echo "insalling test dependencies"
pip2 install -U pytest
pip2 install -r dev_requirements.txt

echo "installing in develop mode to run unit tests"
pip2 install -e .
py.test --cov syncloud test
python setup.py develop --uninstall