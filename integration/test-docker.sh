#!/bin/bash

APP_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
cd ${APP_DIR}

apt-get install docker.io sshpass
service docker start

if [ ! -f rootfs.tar.gz ]; then
  echo "rootfs.tar.gz is not ready, run 'sudo ./bootstrap.sh'"
  exit 1
fi

function sshexec {
    sshpass -p "syncloud" ssh -o StrictHostKeyChecking=no root@localhost -p 2222 "$1"
}

function cleanup {

    echo "cleaning old rootfs"
    rm -rf rootfs

    echo "docker images"
    docker images -q

    echo "removing images"
    docker rm $(docker kill $(docker ps -qa))
    docker rmi $(docker images -q)

    echo "docker images"
    docker images -q
}

cleanup

echo "extracting rootfs"
tar xzf rootfs.tar.gz

#echo "rootfs version: $(<rootfs/version)"

mkdir rootfs/test
rsync -a . rootfs/test --exclude=rootfs*

echo "importing rootfs"
tar -C rootfs -c . | docker import - syncloud

echo "starting rootfs"
docker run --name rootfs --privileged -d -it -p 2222:22 syncloud /sbin/init

sleep 3

echo "running tests"
ssh-keygen -f "/root/.ssh/known_hosts" -R [localhost]:2222

sshexec /test/integration/test.sh