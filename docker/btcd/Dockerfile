FROM hjfreyer/base

# Consolidate
RUN  apt-get install -y golang
RUN  apt-get install -y git
RUN  apt-get install -y mercurial

ENV  GOPATH=/go
RUN  mkdir -p $GOPATH

RUN  go get github.com/btcsuite/btcd/...

RUN  useradd btcserver -m
USER btcserver

CMD  $GOPATH/bin/btcd