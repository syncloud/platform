#!/bin/bash -e

cp /snap/platform/current/certs/* /usr/share/ca-certificates/mozilla/
/usr/sbin/update-ca-certificates