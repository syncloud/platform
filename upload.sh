#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

app=$1
branch=$2
build_number=$3
bucket=$4
s3Key=$5
s3Secret=$6

function upload_file() {

  local file=$1
  local resource="/${bucket}/apps/${file}"
  local contentType="application/x-compressed-tar"
  local dateValue=`date -R`
  local stringToSign="PUT\n\n${contentType}\n${dateValue}\n${resource}"
  local signature=`echo -en ${stringToSign} | openssl sha1 -hmac ${s3Secret} -binary | base64`
  curl -k -X PUT -T "${file}" \
    -H "Host: ${bucket}.s3.amazonaws.com" \
    -H "Date: ${dateValue}" \
    -H "Content-Type: ${contentType}" \
    -H "Authorization: AWS ${s3Key}:${signature}" \
    https://${bucket}.s3.amazonaws.com/apps/${file}
}

mkdir -p /opt/app
SAMCMD=/opt/app/sam/bin/sam

if [ ! -f ${SAMCMD} ]; then
    ${DIR}/install-sam.sh 85 stable
fi

if [ "${branch}" == "master" ] || [ "${branch}" == "stable" ] ; then
   
  upload_file ${app}-${build_number}-x86_64.tar.gz 
  upload_file ${app}-${build_number}-armv7l.tar.gz

  if [ "${branch}" == "stable" ]; then
    branch=rc
  fi

  ${SAMCMD} release $branch $branch --override ${app}=${build_number}

  echo "##teamcity[buildStatus text='released to $branch']"

fi

