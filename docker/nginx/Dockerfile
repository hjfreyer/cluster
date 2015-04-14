FROM debian:jessie

EXPOSE 80
EXPOSE 443

RUN apt-get update && \
    apt-get install -y nginx

RUN useradd -m steve

COPY nginx.conf /etc/nginx/nginx.conf

ADD run.sh /bin/run.sh

CMD ["run.sh"]