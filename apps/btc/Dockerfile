FROM golang:latest

RUN  go get -u github.com/Masterminds/glide

RUN  git clone https://github.com/btcsuite/btcd $GOPATH/src/github.com/btcsuite/btcd && \
 		 git clone https://github.com/btcsuite/btcwallet $GOPATH/src/github.com/btcsuite/btcwallet

RUN  cd $GOPATH/src/github.com/btcsuite/btcd && \
     glide install && \
     go install . ./cmd/...

RUN  cd $GOPATH/src/github.com/btcsuite/btcwallet && \
     glide install && \
     go install . ./cmd/... && \
     go install

ADD start.sh /usr/bin/

ENTRYPOINT ["start.sh"]