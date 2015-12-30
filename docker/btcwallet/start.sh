#!/bin/bash -e

sleep 30
exec /go/bin/btcwallet \
  --datadir=$BTCDIR \
  --username=user \
  --password=$(cat /creds/rpcpass) \
  --rpclisten=:8332 \
  --rpccert=/creds/rpc.cert \
  --rpckey=/creds/rpc.key \
  --mainnet \
  "$@"
