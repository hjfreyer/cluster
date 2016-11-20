FROM golang:latest

EXPOSE 8332

RUN  go get -u github.com/Masterminds/glide

RUN  git clone https://github.com/btcsuite/btcwallet $GOPATH/src/github.com/btcsuite/btcwallet && \
     cd $GOPATH/src/github.com/btcsuite/btcwallet && \
     glide install && \
     go install . ./cmd/... && \
     go install
     
ADD start.sh /usr/bin/

ENTRYPOINT ["start.sh"]
