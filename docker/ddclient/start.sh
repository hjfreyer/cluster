#!/bin/bash

USER=$(cat /creds/username)
PASS=$(cat /creds/password)
HOST=h.hjfreyer.com

while true
do
    curl "https://$USER:$PASS@domains.google.com/nic/update?hostname=$HOST"
    sleep 3600
done

