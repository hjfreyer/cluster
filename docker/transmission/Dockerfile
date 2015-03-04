FROM hjfreyer/base

EXPOSE 9091
EXPOSE 51413

RUN  apt-get install -y transmission-daemon

CMD  ["/usr/bin/transmission-daemon", "-f", "--config-dir=/transmission/config/", "--download-dir=/transmission/data/"]