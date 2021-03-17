#!/bin/sh
set -e

# Copy over files mounted from configmap
cp /configmap/* /conf/

# Test the new config
/usr/sbin/nginx -t -c /conf/nginx.conf

# Signal nginx to gracefully reload
/usr/sbin/nginx -s reload -c /conf/nginx.conf
