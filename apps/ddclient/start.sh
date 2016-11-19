#!/bin/bash

HOST=h.hjfreyer.com

while true
do
    curl "https://$USERNAME:$PASSWORD@domains.google.com/nic/update?hostname=$HOST"
    sleep 3600
done
