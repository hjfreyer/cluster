#!/bin/bash -e

exec /go/bin/btcd \
     -b=/btc/ \
     --rpcuser=user \
     --rpcpass="$(cat /creds/rpcpass)" \
     --notls \
     --rpclisten=127.0.0.1:8334 \
     "$@"
