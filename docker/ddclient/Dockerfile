FROM debian:jessie

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y ddclient

VOLUME /creds

ADD start.sh /

CMD /start.sh 
