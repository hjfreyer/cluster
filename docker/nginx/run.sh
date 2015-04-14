#!/bin/bash -e

sed -i "s/@TRANSMISSION_HOSTPORT@/$TRANSMISSION_HTTP_SERVICE_HOST:$TRANSMISSION_HTTP_SERVICE_PORT/" /etc/nginx/nginx.conf

exec nginx
