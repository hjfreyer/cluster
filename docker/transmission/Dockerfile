FROM hjfreyer/base

RUN  apt-get install -y transmission-daemon

EXPOSE 9091
EXPOSE 51413

RUN  useradd -u 169 transmission

USER transmission

CMD  ["/usr/bin/transmission-daemon", "-f", "--config-dir=/config/", "--download-dir=/data/"]