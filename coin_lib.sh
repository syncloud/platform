#!/usr/bin/env bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

if [ ! -z "$1" ]; then
  ARCH=$1
else
  ARCH=$(dpkg-architecture -qDEB_HOST_GNU_CPU)
fi

if [ ! -d lib ]; then
  mkdir lib
fi

rm -rf lib/*

coin --to=lib py https://pypi.python.org/packages/2.7/r/requests/requests-2.7.0-py2.py3-none-any.whl
coin --to=lib py https://pypi.python.org/packages/py2/u/urllib3/urllib3-1.10.4-py2-none-any.whl
coin --to=lib py https://pypi.python.org/packages/2.7/b/beautifulsoup4/beautifulsoup4-4.4.0-py2-none-any.whl
coin --to=lib py https://pypi.python.org/packages/source/I/IPy/IPy-0.83.tar.gz
coin --to=lib py https://pypi.python.org/packages/source/m/massedit/massedit-0.67.1.zip
coin --to=lib py https://pypi.python.org/packages/source/j/jsonpickle/jsonpickle-0.9.2.tar.gz
coin --to=lib py https://pypi.python.org/packages/source/s/syncloud-lib/syncloud-lib-2.tar.gz

coin --to=lib py https://pypi.python.org/packages/source/p/python-crontab/python-crontab-1.9.3.tar.gz
coin --to=lib py https://pypi.python.org/packages/any/p/python-dateutil/python_dateutil-2.4.2-py2.py3-none-any.whl
coin --to=lib py https://pypi.python.org/packages/3.3/s/six/six-1.9.0-py2.py3-none-any.whl

coin --to=lib py https://pypi.python.org/packages/source/F/Flask/Flask-0.10.1.tar.gz
coin --to=lib py https://pypi.python.org/packages/source/F/Flask-Login/Flask-Login-0.2.11.tar.gz
coin --to=lib py https://pypi.python.org/packages/source/i/itsdangerous/itsdangerous-0.24.tar.gz
coin --to=lib py https://pypi.python.org/packages/2.7/W/Werkzeug/Werkzeug-0.10.4-py2.py3-none-any.whl
coin --to=lib py https://pypi.python.org/packages/source/J/Jinja2/Jinja2-2.7.3.tar.gz

coin --to=lib py http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_MarkupSafe_${ARCH}/lastSuccessful/MarkupSafe-0.23-cp27-none-linux_${ARCH}.whl

coin --to=lib py http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_psutil_${ARCH}/lastSuccessful/psutil-2.1.3-cp27-none-linux_${ARCH}.whl
coin --to=lib py http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_miniupnpc_${ARCH}/lastSuccessful/miniupnpc-1.9-cp27-none-linux_${ARCH}.whl

coin --to=lib py http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_python_ldap_${ARCH}/lastSuccessful/python_ldap-2.4.19-cp27-none-linux_${ARCH}.whl

coin --to=lib raw http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_certbot_${ARCH}/lastSuccessful/certbot-${ARCH}.tar.gz

