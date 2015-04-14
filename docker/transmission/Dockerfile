FROM debian:jessie

EXPOSE 9091
EXPOSE 51413

RUN  apt-get update && \
     apt-get install -y transmission-daemon

RUN  useradd steve -m

USER steve

CMD  ["/usr/bin/transmission-daemon", "-f", "--config-dir=/transmission/config/", "--download-dir=/transmission/data/"]