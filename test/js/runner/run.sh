#!/usr/bin/env bash

APP_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${APP_DIR}

sudo apt-get install phantomjs

phantomjs QUnitTeamCityDriver.phantom.js qunit.html