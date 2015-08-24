FROM debian:jessie

EXPOSE 8333
EXPOSE 8334

ENV  BTCDIR /btc
ENV  GOPATH /go

RUN apt-get update && apt-get install -y \
    git \
    golang \
    supervisor

RUN mkdir /go
RUN mkdir -p /var/log/supervisor

RUN  go get github.com/btcsuite/btcd/...

COPY supervisord.conf /etc/supervisord.conf

ENTRYPOINT ["/go/bin/btcd", "-b=$BTCDIR", "--rpcuser=user", "--rpcpass=pass", "--rpclisten=:8334", "--debuglevel=debug"]
#CMD ["/usr/bin/supervisord"]
