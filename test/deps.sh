#!/bin/bash -e

apt-get update
apt-get install -y sshpass openssh-client netcat-openbsd file curl apache2-utils
pip install -r requirements.txt
