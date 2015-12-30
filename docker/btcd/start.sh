#!/bin/bash -e

exec /go/bin/btcd \
  -b=$BTCDIR \
  --rpcuser=user \
  --rpcpass=$(cat /creds/rpcpass) \
  --rpccert=/creds/rpc.cert \
  --rpckey=/creds/rpc.key \
  "$@"
