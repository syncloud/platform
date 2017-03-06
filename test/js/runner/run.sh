#!/usr/bin/env bash

APP_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${APP_DIR}

xvfb-run --server-args="-screen 0, 1024x4096x24" phantomjs QUnitTeamCityDriver.phantom.js qunit.html