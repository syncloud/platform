#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}
NAME="platform"

if [[ -z "$1" ]]; then
    echo "usage $0 version"
    exit 1
fi

ARCH=$(dpkg-architecture -q DEB_HOST_GNU_CPU)
#temp fix until we move to cpu arch
if [ "${ARCH}" == 'arm' ]; then
    ARCH="armv7l"
fi
VERSION=$1

cd ${DIR}

BUILD_DIR=${DIR}/build/${NAME}
SNAP_DIR=${DIR}/build/snap
rm -rf build
mkdir -p ${BUILD_DIR}

DOWNLOAD_URL=http://build.syncloud.org:8111/guestAuth/repository/download
coin --to ${BUILD_DIR} raw ${DOWNLOAD_URL}/thirdparty_nginx_${ARCH}/lastSuccessful/nginx-${ARCH}.tar.gz
coin --to ${BUILD_DIR} raw ${DOWNLOAD_URL}/thirdparty_uwsgi_${ARCH}/lastSuccessful/uwsgi-${ARCH}.tar.gz
coin --to ${BUILD_DIR} raw ${DOWNLOAD_URL}/thirdparty_openldap_${ARCH}/lastSuccessful/openldap-${ARCH}.tar.gz
coin --to=${BUILD_DIR} raw http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_openssl_${ARCH}/lastSuccessful/openssl-${ARCH}.tar.gz
coin --to ${BUILD_DIR} raw ${DOWNLOAD_URL}/thirdparty_python_${ARCH}/lastSuccessful/python-${ARCH}.tar.gz

${BUILD_DIR}/python/bin/pip install -r ${DIR}/requirements.txt

cd src
rm -f version
echo ${VERSION} >> version
${BUILD_DIR}/python/bin/python setup.py install
cd ..

cp -r ${DIR}/bin ${BUILD_DIR}
cp -r ${DIR}/config ${BUILD_DIR}/config.templates
cp -r ${DIR}/www ${BUILD_DIR}

mkdir ${BUILD_DIR}/META
echo ${NAME} >> ${BUILD_DIR}/META/app
echo ${VERSION} >> ${BUILD_DIR}/META/version

echo "zipping"
rm -rf ${NAME}*.tar.gz
tar cpzf ${DIR}/${NAME}-${VERSION}-${ARCH}.tar.gz -C ${DIR}/build/ ${NAME}

echo "snapping"
ARCH=$(dpkg-architecture -q DEB_HOST_ARCH)
rm -rf ${DIR}/*.snap
mkdir ${SNAP_DIR}
cp -r ${BUILD_DIR}/* ${SNAP_DIR}/
sed -i 's/installer:.*/installer: snapd/g' ${SNAP_DIR}/config.templates/platform.cfg
cp -r ${DIR}/snap/meta ${SNAP_DIR}/
cp ${DIR}/snap/snap.yaml ${SNAP_DIR}/meta/snap.yaml
echo "version: $VERSION" >> ${SNAP_DIR}/meta/snap.yaml
echo "architectures:" >> ${SNAP_DIR}/meta/snap.yaml
echo "- ${ARCH}" >> ${SNAP_DIR}/meta/snap.yaml

mksquashfs ${SNAP_DIR} ${DIR}/platform_${VERSION}_${ARCH}.snap -noappend -comp xz -no-xattrs -all-root

${DIR}/unit-test.sh



