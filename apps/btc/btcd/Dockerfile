FROM golang:latest

RUN  go get -u github.com/Masterminds/glide

RUN  git clone https://github.com/btcsuite/btcd $GOPATH/src/github.com/btcsuite/btcd && \
     cd $GOPATH/src/github.com/btcsuite/btcd && \
     glide install && \
     go install . ./cmd/...

ADD start.sh /usr/bin/

ENTRYPOINT ["start.sh"]
