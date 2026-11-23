#!/bin/bash -e

SHIPPED=/snap/platform/current/certs
for f in ${SHIPPED}/*.crt; do
  name=$(basename $f .crt)
  [ -e /etc/ssl/certs/${name}.pem ] || echo "not enabled: ${name}"
done
[ "$(ls ${SHIPPED})" = "$(ls /usr/share/ca-certificates/mozilla)" ] || echo "mozilla dir differs from ${SHIPPED}"
