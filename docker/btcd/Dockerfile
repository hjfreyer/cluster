FROM golang

EXPOSE 8333
EXPOSE 8334

ENV  BTCDIR=/btc

RUN  go get github.com/btcsuite/btcd/...

RUN  useradd box -m
USER box

CMD ["btcd", "-b=$BTCDIR", "--rpcuser=user", "--rpcpass=pass", "--rpclisten=:8334"]