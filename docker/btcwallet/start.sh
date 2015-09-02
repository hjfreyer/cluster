#!/bin/bash -e

exec /go/bin/btcwallet \
  --datadir=$BTCDIR \
  --rpcuser=user \
  --rpcpass=$(cat /creds/rpcpass) \
  --rpclisten=:8332 \
  --rpccert=/creds/rpc.cert \
  --rpckey=/creds/rpc.key \
  --mainnet \
  "$@"
