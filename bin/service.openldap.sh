#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

$DIR/openldap/libexec/slapd -h ldap://127.0.0.1:389 -F ${SNAP_COMMON}/slapd.d


