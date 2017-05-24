#!/bin/bash -xe

apt-get update
apt-get install -qq squashfs-tools dpkg-dev python-dev libsasl2-dev libldap2-dev libssl-dev libffi-dev apache2-utils wget unzip
wget https://bootstrap.pypa.io/get-pip.py
python get-pip.py
pip install coin
ARCH=$(uname -m)
if [ $ARCH == "x86_64" ]; then
  wget --progress dot:giga https://bitbucket.org/ariya/phantomjs/downloads/phantomjs-2.1.1-linux-x86_64.tar.bz2
  tar xjf phantomjs-2.1.1-linux-x86_64.tar.bz2
  cp ./phantomjs-2.1.1-linux-x86_64/bin/phantomjs /usr/bin
else
  wget --progress dot:giga https://github.com/fg2it/phantomjs-on-raspberry/releases/download/v2.1.1-wheezy-jessie-armv6/phantomjs
  cp phantomjs /usr/bin
fi
chmod +x /usr/bin/phantomjs