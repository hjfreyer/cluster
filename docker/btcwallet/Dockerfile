FROM golang:latest

EXPOSE 8332

ENV  BTCDIR /btc

RUN  go get github.com/btcsuite/btcwallet/...

ADD start.sh /usr/bin/

ENTRYPOINT ["start.sh"]
