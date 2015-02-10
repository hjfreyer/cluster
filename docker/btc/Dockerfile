FROM hjfreyer/base

RUN  apt-get install -y bitcoind

EXPOSE 8333

RUN  useradd -u 169 transmission

USER transmission

CMD  ["/usr/bin/transmission-daemon", "-f", "--config-dir=/config/", "--download-dir=/data/"]