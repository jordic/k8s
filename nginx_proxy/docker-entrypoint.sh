#!/bin/bash
set -e

# build configuration for available services
/main -path=/etc/nginx/conf.d

if [ "$1" = 'nginx-daemon' ]; then
    exec nginx -g "daemon off;";
fi

exec "$@"
