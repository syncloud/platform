#!/usr/bin/env bash


APP_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

if [ ! -f ${APP_DIR}/rootfs.tar.gz ]; then
  echo "rootfs.tar.gz is not ready, run 'sudo ./bootstrap.sh'"
  exit 1
fi

apt-get install docker.io
service docker start

function cleanup {

    echo "cleaning old rootfs"
    rm -rf /tmp/rootfs

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
tar xzf ${APP_DIR}/rootfs.tar.gz -C /tmp

#echo "rootfs version: $(<rootfs/version)"
sed -i 's/Port 22/Port 2222/g' /tmp/rootfs/etc/ssh/sshd_config
mkdir /tmp/rootfs/test

echo "copying all files to rootfs"
rsync -a ${APP_DIR}/ /tmp/rootfs/test --exclude=/rootfs* --exclude=/dist --exclude=/build --exclude=/nginx --exclude=/uwsgi

echo "importing rootfs"
tar -C /tmp/rootfs -c . | docker import - syncloud

echo "starting rootfs"
docker run --net host -v /var/run/dbus:/var/run/dbus --name rootfs --privileged -d -it syncloud /sbin/init

echo "sleeping for services to start"
sleep 5

ssh-keygen -f "/root/.ssh/known_hosts" -R [localhost]:2222