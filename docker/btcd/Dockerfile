FROM golang

ENV  BTCDIR=/btc

RUN  go get github.com/btcsuite/btcd/...

RUN  useradd box -m
USER box

ENTRYPOINT ["btcd", "-b=$BTCDIR"]