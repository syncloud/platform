#!/bin/bash

sudo apt-get install -y snapcraft ruby
sudo gem install jekyll

ARCH=$(dpkg-architecture -q DEB_HOST_ARCH)
if [ $ARCH == "amd64" ]; then
  wget https://github.com/fg2it/tarphantomjs-on-raspberry/releases/download/v2.1.1-wheezy-jessie-armv6/phantomjs
  cp phantomjs /usr/bin
else
  wget https://bitbucket.org/ariya/phantomjs/downloads/phantomjs-2.1.1-linux-x86_64.tar.bz2
  tar xjvf https://bitbucket.org/ariya/phantomjs/downloads/phantomjs-2.1.1-linux-x86_64.tar.bz2
  cp ./phantomjs-2.1.1-linux-x86_64/bin/phantomjs /usr/bin

fi