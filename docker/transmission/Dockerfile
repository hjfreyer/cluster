FROM debian:jessie

RUN apt-get update

EXPOSE 9091
EXPOSE 51413

RUN  apt-get install -y transmission-daemon

RUN  useradd steve -m

USER steve

CMD  ["/usr/bin/transmission-daemon", "-f", "--config-dir=/transmission/config/", "--download-dir=/transmission/data/"]