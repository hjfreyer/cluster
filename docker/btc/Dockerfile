FROM debian:sid

EXPOSE 8332
EXPOSE 8333

RUN apt-get update && apt-get install -y bitcoind

VOLUME /creds

ADD start.sh /bin/

CMD start.sh 
