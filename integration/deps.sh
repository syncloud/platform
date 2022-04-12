#!/bin/bash -e

apt-get update
apt-get install -y sshpass openssh-client netcat rustc file libxml2-dev libxslt-dev build-essential libz-dev curl
pip install -r requirements.txt
