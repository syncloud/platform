#!/bin/bash -xe

apt-get update
apt-get install -qq squashfs-tools dpkg-dev python-dev libsasl2-dev libldap2-dev libssl-dev libffi-dev apache2-utils wget unzip sshpass xvfb curl netcat libfontconfig libgtk-3-0 libasound2
wget https://bootstrap.pypa.io/get-pip.py
python get-pip.py
pip install coin
ARCH=$(uname -m)
if [ $ARCH == "x86_64" ]; then
  wget --progress dot:giga http://artifact.syncloud.org/3rdparty/phantomjs-2.1.1-linux-x86_64.tar.bz2
  tar xjf phantomjs-2.1.1-linux-x86_64.tar.bz2
  cp ./phantomjs-2.1.1-linux-x86_64/bin/phantomjs /usr/bin
else
  wget --progress dot:giga http://artifact.syncloud.org/3rdparty/phantomjs-2.1.1-armhf
  cp phantomjs-2.1.1-armhf /usr/bin/phantomjs
fi
chmod +x /usr/bin/phantomjs