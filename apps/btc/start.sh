#!/bin/bash -e

CMD=$1
shift

case $CMD in
		btcd)
				exec /go/bin/btcd \
						 -b=/btc/ \
						 --rpcuser=user \
						 --rpcpass=$(cat /creds/rpcpass) \
						 --notls \
						 --rpclisten=127.0.0.1:8334 \
						 "$@"
				;;

		btcwallet)
				WALLET_FILE=/wallet/mainnet/wallet.db

				if [[ ! -f "$WALLET_FILE" ]]; then
						mkdir -p "$(dirname "$WALLET_FILE")"
						cp /source/wallet.db "$WALLET_FILE"
				fi

				exec /go/bin/btcwallet \
						 -A /wallet/ \
						 --username=user \
						 --password=$(cat /creds/rpcpass) \
						 --btcdusername=user \
						 --btcdpassword=$(cat /creds/rpcpass) \
						 --noclienttls \
						 --rpclisten=:8332 \
						 --rpccert=/creds/rpc.cert \
						 --rpckey=/creds/rpc.key \
						 "$@"
				;;

		*)
				echo "Unknown command: $CMD"
				exit 1
				;;
esac

