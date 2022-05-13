FROM golang:1.16-alpine as backbuilder

ARG STAGE_MODE=prod

WORKDIR /app
COPY ./src/server/go.mod ./src/server/go.sum ./
RUN go mod download

COPY ./src/server .
RUN mkdir -p ./build &&\
    go build -o ./build/server -v -ldflags="-s -w" ./cmd/main.go

RUN mkdir -p ./build/misc &&\
    cp ./misc/sitemap.json ./build/misc/sitemap.json

FROM node:16.13-alpine as frontbuilder

ARG STAGE_MODE=prod

WORKDIR /app
RUN mkdir -p ./build
COPY ./src/server/views ./build/views
COPY ./src/server/views ./src/server/views

WORKDIR /app/src/frontend
COPY ./src/frontend/package.json ./src/frontend/package-lock.json ./
RUN npm install

ARG NODE_ENV=production

COPY ./$STAGE_MODE.env ./.env
COPY ./src/frontend .
RUN NODE_ENV=$NODE_ENV npm run $STAGE_MODE
RUN rm .env

FROM alpine:3

ARG STAGE_MODE=prod
ENV TZ=Europe/Moscow

RUN apk update &&\
    apk --no-cache add ca-certificates tzdata &&\
    update-ca-certificates

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime &&\
    echo $TZ > /etc/timezone

COPY --from=backbuilder ./app/build /home/app
COPY --from=frontbuilder ./app/build /home/app
COPY ./data/certbot /etc/letsencrypt

RUN mkdir -p ./home/app/logs
