#!/bin/bash -e
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
/bin/rm -f /var/snap/testapp/common/web.socket
exec ${DIR}/nginx/bin/nginx.sh -c /var/snap/testapp/current/config/nginx.conf -p ${DIR}/nginx -e stderr
