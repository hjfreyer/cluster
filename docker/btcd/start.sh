#!/bin/bash -e

exec /go/bin/btcd \
  -b=$BTCDIR \
  --rpcuser=user \
  --rpcpass=$(cat /creds/rpcpass) \
  --rpclisten=:8334 \
  --rpccert=/creds/rpc.cert \
  --rpckey=/creds/rpc.key \
  "$@"
