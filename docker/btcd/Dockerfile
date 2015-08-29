FROM golang:latest

EXPOSE 8333
EXPOSE 8334

ENV  BTCDIR /btc

RUN  go get github.com/btcsuite/btcd/...

ADD start.sh /usr/bin/

ENTRYPOINT ["start.sh"]
