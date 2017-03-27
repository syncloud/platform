#!/bin/bash -xe

APP_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${APP_DIR}

phantomjs QUnitTeamCityDriver.phantom.js qunit.html
