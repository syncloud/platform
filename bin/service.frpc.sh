#!/bin/sh

CONFIG=${SNAP_DATA}/frpc.toml
if [ -f "${CONFIG}" ]; then
    exec ${SNAP}/frp/frpc -c "${CONFIG}"
fi
echo "relay tunnel disabled, frpc idle (no ${CONFIG})"
exec sleep infinity
