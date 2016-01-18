FROM ubuntu:14.04
MAINTAINER yt@hjfreyer.com

EXPOSE 22

RUN useradd hjfreyer \
  -u 1000 \
  -s /bin/bash

RUN apt-get update && \
  apt-get install -y \
  openssh-server \
  apache2 \
  supervisor \
  docker.io

RUN echo "hjfreyer:foo" | chpasswd
RUN mkdir -p /var/run/sshd

ADD supervisord.conf /etc/supervisor/conf.d/supervisord.conf

ENTRYPOINT supervisord -c /etc/supervisor/conf.d/supervisord.conf
