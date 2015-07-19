FROM golang

EXPOSE 8333
EXPOSE 8334

ENV  BTCDIR /btc

RUN  go get github.com/btcsuite/btcd/...


RUN  useradd box -m

RUN  mkdir $BTCDIR && chmod 777 $BTCDIR

USER box

CMD ["btcd", "-b=$BTCDIR", "--rpcuser=user", "--rpcpass=pass", "--rpclisten=:8334"]