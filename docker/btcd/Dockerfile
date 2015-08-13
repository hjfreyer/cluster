FROM debian:jessie

EXPOSE 8333
EXPOSE 8334

ENV  BTCDIR /btc

RUN apt-get update
RUN apt-get install -y golang
RUN apt-get install -y git

RUN mkdir /go

ENV GOPATH /go

RUN  go get github.com/btcsuite/btcd/...

CMD ["btcd", "-b=$BTCDIR", "--rpcuser=user", "--rpcpass=pass", "--rpclisten=:8334"]