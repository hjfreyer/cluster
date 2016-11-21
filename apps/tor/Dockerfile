FROM debian:jessie

RUN  apt-get update && apt-get install -y tor

ADD  torrc /etc/tor/torrc

USER debian-tor

CMD  tor