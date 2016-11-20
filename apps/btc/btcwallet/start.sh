#!/bin/bash -e

WALLET_FILE=/wallet/mainnet/wallet.db

if [[ ! -f "$WALLET_FILE" ]]; then
    mkdir -p "$(dirname "$WALLET_FILE")"
    cp /source/wallet.db "$WALLET_FILE"
fi

exec /go/bin/btcwallet \
     -A /wallet/ \
     --username=user \
     --password=$(cat /creds/rpcpass) \
     --noclienttls \
     --rpclisten=:8332 \
     --rpccert=/creds/rpc.cert \
     --rpckey=/creds/rpc.key \
     "$@"
