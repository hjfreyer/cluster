FROM hjfreyer/base

EXPOSE 80
EXPOSE 443

RUN apt-get install -y nginx

COPY nginx.conf /etc/nginx/nginx.conf

ADD run.sh /bin/run.sh

CMD ["run.sh"]