#!/bin/sh
set -e

# Copy over files mounted from configmap
cp /configmap/* /conf/

/watcher &
exec /usr/sbin/nginx -c /conf/nginx.conf
