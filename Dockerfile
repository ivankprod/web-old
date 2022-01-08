FROM ubuntu

ARG STAGE_MODE=prod
ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Europe/Moscow

RUN apt-get update
RUN apt-get install -y ca-certificates tzdata
RUN update-ca-certificates

RUN ln -fs /usr/share/zoneinfo/Europe/Moscow /etc/localtime
RUN dpkg-reconfigure --frontend noninteractive tzdata

COPY ./build_$STAGE_MODE /home/app
COPY ./data/certbot /etc/letsencrypt
