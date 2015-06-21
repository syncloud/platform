#!/bin/sh -x

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

NAME="platform"
ARCHITECTURE=$1
VERSION="local"
if [ ! -z "$2" ]; then
    VERSION=$2
fi

if ! jekyll -v; then
  echo "installing jekyll"
  apt-get -y install ruby ruby-dev make gcc nodejs
  gem install jekyll --no-rdoc --no-ri
fi

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

wget -O python.tar.gz http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_python_${ARCHITECTURE}/lastSuccessful/python.tar.gz
tar -xvf python.tar.gz
rm python.tar.gz

PYTHON_PATH='python/bin'

wget -O get-pip.py https://bootstrap.pypa.io/get-pip.py
${PYTHON_PATH}/python get-pip.py
rm get-pip.py

${PYTHON_PATH}/pip install wheel

PSUTIL_WHL="psutil-2.1.3-cp27-none-linux_${ARCHITECTURE}.whl"
wget http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_psutil_${ARCHITECTURE}/lastSuccessful/${PSUTIL_WHL}
${PYTHON_PATH}/pip install ${PSUTIL_WHL}
rm ${PSUTIL_WHL}

PYTHON_LDAP_WHL="python_ldap-2.4.19-cp27-none-linux_${ARCHITECTURE}.whl"
wget http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_python_ldap_${ARCHITECTURE}/lastSuccessful/${PYTHON_LDAP_WHL}
${PYTHON_PATH}/pip install ${PYTHON_LDAP_WHL}
rm ${PYTHON_LDAP_WHL}

${PYTHON_PATH}/pip install ${DIR}/src/dist/syncloud-platform-${VERSION}.tar.gz

wget -O nginx.tar.gz http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_nginx_${ARCHITECTURE}/lastSuccessful/nginx.tar.gz
tar -xvf nginx.tar.gz
rm nginx.tar.gz

wget -O uwsgi.tar.gz http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_uwsgi_${ARCHITECTURE}/lastSuccessful/uwsgi.tar.gz
tar -xvf uwsgi.tar.gz
rm uwsgi.tar.gz

wget -O openldap.tar.gz http://build.syncloud.org:8111/guestAuth/repository/download/thirdparty_openldap_${ARCHITECTURE}/lastSuccessful/openldap.tar.gz
tar -xvf openldap.tar.gz
rm openldap.tar.gz

cd ../..

cp -r bin build/${NAME}
cp -r config build/${NAME}
cp -r www build/${NAME}

mkdir build/${NAME}/META
echo ${NAME} >> build/${NAME}/META/app
echo ${VERSION} >> build/${NAME}/META/version

cd build
tar -zcvf ${NAME}-${VERSION}-${ARCHITECTURE}.tar.gz ${NAME}
