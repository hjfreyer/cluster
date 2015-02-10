#!/bin/bash -e

NAME=transmission
CONFIG=/var/lib/transmission/.config/transmission-daemon/
DATA=/home/download/torrent/

docker rm -f $NAME
docker build -t $NAME .

docker run -d \
    --name $NAME \
    -v $CONFIG:/config/ -v $DATA:/data/ \
    -p 127.0.0.1:9091:9091 -p 51413:51413 \
    transmission
