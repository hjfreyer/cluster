#!/bin/bash -e

echo "
protocol=dyndns2
use=web
server=domains.google.com
ssl=yes
login=$(cat /creds/username)
password=$(cat /creds/password)
h.hjfreyer.com" > /etc/ddclient.conf

exec ddclient --daemon 300 --foreground
