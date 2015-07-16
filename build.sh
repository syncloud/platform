#!/bin/bash -x

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}
NAME="platform"

ARCHITECTURE=$(dpkg-architecture -qDEB_HOST_GNU_CPU)
if [ ! -z "$1" ]; then
    ARCHITECTURE=$1
fi

VERSION="local"
if [ ! -z "$2" ]; then
    VERSION=$2
fi

if ! jekyll -v; then
  echo "installing jekyll"
  apt-get -y install ruby ruby-dev make gcc nodejs
  gem install jekyll --no-rdoc --no-ri
fi

function 3rdparty {
  APP_ID=$1
  APP_FILE=$2
  if [ ! -d ${DIR}/3rdparty ]; then
    mkdir ${DIR}/3rdparty
  fi
  if [ ! -f ${DIR}/3rdparty/${APP_FILE} ]; then
    wget http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_${APP_ID}_${ARCHITECTURE}/lastSuccessful/${APP_FILE} \
    -O ${DIR}/3rdparty/${APP_FILE} --progress dot:giga
  else
    echo "skipping ${APP_ID}"
  fi
}

PSUTIL_WHL="psutil-2.1.3-cp27-none-linux_${ARCHITECTURE}.whl"
PYTHON_LDAP_WHL="python_ldap-2.4.19-cp27-none-linux_${ARCHITECTURE}.whl"
MINIUPNPC_WHL="miniupnpc-1.9-cp27-none-linux_${ARCHITECTURE}.whl"
NGINX_ZIP=nginx.tar.gz
UWSGI_ZIP=uwsgi.tar.gz
OPENLDAP_ZIP=openldap.tar.gz
PYTHON_ZIP=python.tar.gz

3rdparty nginx ${NGINX_ZIP}
3rdparty uwsgi ${UWSGI_ZIP}
3rdparty openldap ${OPENLDAP_ZIP}
3rdparty python ${PYTHON_ZIP}
3rdparty psutil ${PSUTIL_WHL}
3rdparty miniupnpc ${MINIUPNPC_WHL}
3rdparty python_ldap ${PYTHON_LDAP_WHL}

cd www
rm -rf _site
jekyll build
cd ..

rm -f src/version
echo ${VERSION} >> src/version
cd src
python setup.py sdist
cd ..

rm -rf build
mkdir build
mkdir build/${NAME}
cd build/${NAME}

tar -xf ${DIR}/3rdparty/${PYTHON_ZIP}
PYTHON_PATH='python/bin'

wget -O get-pip.py https://bootstrap.pypa.io/get-pip.py
${PYTHON_PATH}/python get-pip.py
rm get-pip.py

${PYTHON_PATH}/pip install wheel
${PYTHON_PATH}/pip install ${DIR}/3rdparty/${PSUTIL_WHL}
${PYTHON_PATH}/pip install ${DIR}/3rdparty/${PYTHON_LDAP_WHL}
${PYTHON_PATH}/pip install ${DIR}/src/dist/syncloud-platform-${VERSION}.tar.gz

tar -xzf ${DIR}/3rdparty/${NGINX_ZIP}
tar -xzf ${DIR}/3rdparty/${UWSGI_ZIP}
tar -xzf ${DIR}/3rdparty/${OPENLDAP_ZIP}

cd ../..

cp -r bin build/${NAME}
cp -r config build/${NAME}
cp -r www build/${NAME}

mkdir build/${NAME}/META
echo ${NAME} >> build/${NAME}/META/app
echo ${VERSION} >> build/${NAME}/META/version

echo "zipping"
tar cpzf ${NAME}-${VERSION}-${ARCHITECTURE}.tar.gz -C build/ ${NAME}