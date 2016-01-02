#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

export TMPDIR=/tmp
export TMP=/tmp
export DEBIAN_FRONTEND=noninteractive

if [ "$#" -ne 7 ]; then
    echo "usage $0 redirect_user redirect_password redirect_domain app_version app_arch sam_version release"
    exit 1
fi

./docker.sh $6

apt-get install -y sshpass
pip2 install -r ${DIR}/../src/dev_requirements.txt
py.test -s verify.py --email=$1 --password=$2 --domain=$3 --app-version=$4 --arch=$5 --release=$7