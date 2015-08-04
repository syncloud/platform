#!/usr/bin/env bash

ROOTFS=/tmp/platform/rootfs
APP_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
cd ${APP_DIR}
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

ARCH=$(dpkg-architecture -q DEB_HOST_GNU_CPU)
SAM_VERSION=$1

SAM=sam-${SAM_VERSION}-${ARCH}.tar.gz

if [ ! -d 3rdparty ]; then
  mkdir 3rdparty
fi

cd 3rdparty
if [ ! -f rootfs-${ARCH}.tar.gz ]; then
  wget http://build.syncloud.org:8111/guestAuth/repository/download/debian_rootfs_${ARCH}/lastSuccessful/rootfs.tar.gz\
  -O rootfs-${ARCH}.tar.gz --progress dot:giga
else
  echo "skipping rootfs"
fi
if [ ! -f ${SAM} ]; then
  wget http://apps.syncloud.org/apps/${SAM} --progress=dot:giga
else
  echo "skipping sam"
fi
cd ..

apt-get install docker.io
service docker start

function cleanup {

    mount | grep rootfs
    mount | grep rootfs | awk '{print "umounting "$1; system("umount "$3)}'
    mount | grep rootfs

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
rm -rf ${ROOTFS}
mkdir -p ${ROOTFS}
tar xzf ${APP_DIR}/3rdparty/rootfs-${ARCH}.tar.gz -C ${ROOTFS}

tar xzf ${APP_DIR}/3rdparty/${SAM} -C ${ROOTFS}/opt/app

sed -i 's/Port 22/Port 2222/g' ${ROOTFS}/etc/ssh/sshd_config

echo "importing rootfs"
tar -C ${ROOTFS} -c . | docker import - syncloud

echo "starting rootfs"
docker run --net host -v /var/run/dbus:/var/run/dbus --name rootfs --privileged -d -it syncloud /sbin/init

echo "sleeping for services to start"
sleep 10

ssh-keygen -f "/root/.ssh/known_hosts" -R [localhost]:2222