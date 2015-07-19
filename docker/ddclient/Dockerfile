FROM debian:jessie

RUN apt-get update && apt-get install -y curl

VOLUME /creds

ADD start.sh /bin/

CMD start.sh 
