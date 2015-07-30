#!/bin/bash -x

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}
NAME="platform"

if [[ -z "$1" || -z "$2" ]]; then
    echo "usage $0 app_arch app_version"
    exit 1
fi

ARCH=$1
VERSION=$2

function 3rdparty {
  APP_ID=$1
  APP_FILE=$2
  if [ ! -d ${DIR}/3rdparty ]; then
    mkdir ${DIR}/3rdparty
  fi
  if [ ! -f ${DIR}/3rdparty/${APP_FILE} ]; then
    wget http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_${APP_ID}_${ARCH}/lastSuccessful/${APP_FILE} \
    -O ${DIR}/3rdparty/${APP_FILE} --progress dot:giga
  else
    echo "skipping ${APP_ID}"
  fi
}

NGINX_ZIP=nginx.tar.gz
UWSGI_ZIP=uwsgi.tar.gz
OPENLDAP_ZIP=openldap.tar.gz
PYTHON_ZIP=python.tar.gz
#JEKYLL_ZIP=jekyll.tar.gz

3rdparty nginx ${NGINX_ZIP}
3rdparty uwsgi ${UWSGI_ZIP}
3rdparty openldap ${OPENLDAP_ZIP}
3rdparty python ${PYTHON_ZIP}
#3rdparty jekyll ${JEKYLL_ZIP}

tar xzf ${DIR}/3rdparty/${JEKYLL_ZIP} -C ${DIR}/3rdparty/
cd www
rm -rf _site
jekyll build
#${DIR}/3rdparty/jekyll/bin/jekyll build
cd ..

wget -O get-pip.py https://bootstrap.pypa.io/get-pip.py
python get-pip.py
rm get-pip.py

rm -f src/version
echo ${VERSION} >> src/version

cd src
python setup.py sdist
cd ..

pip install --upgrade coin
./coin_lib.sh ${ARCH}
coin  --to ${DIR}/lib py ${DIR}/src/dist/syncloud-platform-${VERSION}.tar.gz

BUILD_DIR=${DIR}/build/${NAME}
rm -rf build
mkdir -p ${BUILD_DIR}

tar -xf ${DIR}/3rdparty/${PYTHON_ZIP} -C ${BUILD_DIR}
tar -xzf ${DIR}/3rdparty/${NGINX_ZIP} -C ${BUILD_DIR}
tar -xzf ${DIR}/3rdparty/${UWSGI_ZIP} -C ${BUILD_DIR}
tar -xzf ${DIR}/3rdparty/${OPENLDAP_ZIP} -C ${BUILD_DIR}

cp -r ${DIR}/bin ${BUILD_DIR}
cp -r ${DIR}/config ${BUILD_DIR}
cp -r ${DIR}/www ${BUILD_DIR}
cp -r ${DIR}/lib ${BUILD_DIR}

mkdir ${BUILD_DIR}/META
echo ${NAME} >> ${BUILD_DIR}/META/app
echo ${VERSION} >> ${BUILD_DIR}/META/version

echo "zipping"
rm -rf ${NAME}*.tar.gz
tar cpzf ${DIR}/${NAME}-${VERSION}-${ARCH}.tar.gz -C ${DIR}/build/ ${NAME}