#!/bin/bash -e

wget https://github.com/s3tools/s3cmd/archive/master.zip
unzip master.zip

apt-get install -y python-dateutil python-magic
rm -rf /opt/s3cmd-master
rm -rf /usr/bin/s3cmd
mv s3cmd-master /opt/s3cmd-master
ln -s /opt/s3cmd-master/s3cmd /usr/bin/s3cmd