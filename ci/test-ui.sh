#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )

MODE=$1
DOMAIN=$2

if [[ -z "${MODE}" || -z "${DOMAIN}" ]]; then
    echo "usage $0 mode domain"
    exit 1
fi

for i in 1 2 3 4 5; do
    apt-get -o Acquire::Retries=10 update && \
        apt-get -o Acquire::Retries=10 install -y sshpass sudo && break
    echo "apt attempt ${i} failed, retrying in 60s"
    sleep 60
done
command -v sshpass

for alias in \
    "auth.${DOMAIN}.redirect" \
    "${DOMAIN}.redirect" \
    "unknown.${DOMAIN}.redirect" \
    "externalapp.${DOMAIN}.redirect" \
    "testapp.${DOMAIN}.redirect" \
    "files.${DOMAIN}.redirect"; do
    getent hosts "${DOMAIN}" | sed "s/${DOMAIN}/${alias}/g" >> /etc/hosts
done
cat /etc/hosts

cd "${DIR}/web/e2e"
npm ci
npx playwright test --project="${MODE}"
