#!/bin/bash -e

attempts=100
attempt=0
CMD="sshpass -p syncloud ssh -o StrictHostKeyChecking=no root@device"

set +e
${CMD} date
while test $? -gt 0
do
  if [[ ${attempt} -gt ${attempts} ]]; then
    exit 1
  fi
  sleep 3
  echo "Waiting for SSH $attempt"
  attempt=$((attempt+1))
  ${CMD} date
done
set -e
