#!/usr/bin/env bash

APP_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
ROOTFS=${APP_DIR}/rootfs
cd ${APP_DIR}
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

ARCH=$(dpkg-architecture -q DEB_HOST_GNU_CPU)
SAM_VERSION=$1
RELEASE=$2

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

apt-get install -y docker.io apache2-utils
service docker start

function cleanup {

    losetup -a
#    mount
    losetup -d /dev/loop0
    losetup -a
    mount | grep rootfs | awk '{print "umounting "$1; system("umount "$3)}'
#    mount | grep docker | awk '{print "umounting "$1; system("umount "$3)}'
    mount | grep rootfs

    echo "cleaning old rootfs"
    rm -rf ${ROOTFS}

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

cp -r ${APP_DIR}/integration ${ROOTFS}

echo "importing rootfs"
tar -C ${ROOTFS} -c . | docker import - syncloud

echo "starting rootfs"
docker run --net host -v /var/run/dbus:/var/run/dbus --name rootfs --privileged -d -it syncloud /sbin/init --cap-add=ALL

sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost date
while test $? -gt 0
do
  sleep 1
  echo "Waiting for SSH ..."
  sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost date
done

sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost "/opt/app/sam/bin/sam update --release ${RELEASE}"

ssh-keygen -f "/root/.ssh/known_hosts" -R [localhost]:2222
