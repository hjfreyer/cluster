#!/bin/bash -e

echo $TRANSMISSION_HTTP_SERVICE_HOST
echo $TRANSMISSION_HTTP_SERVICE_PORT

sed "s/@TRANSMISSION_HOSTPORT@/$TRANSMISSION_HTTP_SERVICE_HOST:$TRANSMISSION_HTTP_SERVICE_PORT/" /etc/nginx/nginx.conf > /tmp/conf
mv /tmp/conf /etc/nginx/nginx.conf

exec nginx
