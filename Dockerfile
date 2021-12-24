FROM ubuntu

ARG STAGE_MODE=prod

RUN apt-get update
RUN apt-get install -y ca-certificates
RUN update-ca-certificates

COPY ./build_$STAGE_MODE /home/app
COPY ./data/certbot /etc/letsencrypt
