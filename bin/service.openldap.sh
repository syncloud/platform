#!/bin/bash


DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 [action]"
    exit 1
fi

export LD_LIBRARY_PATH=$DIR/openldap/lib
export SASL_PATH=$DIR/openldap/lib
export SASL_CONF_PATH=$DIR/openldap/lib

SOCKET="${SNAP_DATA}/openldap.socket"
case $1 in
start)
    exec ${DIR}/openldap/sbin/slapd.sh -h "ldap://127.0.0.1:389 ldapi://${SOCKET//\//%2F}" -F /var/snap/platform/current/slapd.d
    ;;
*)
    echo "not valid command"
    exit 1
    ;;
esac
