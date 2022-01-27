FROM golang:1.16-alpine as backbuilder

ARG STAGE_MODE=prod

RUN apk update

WORKDIR /app
COPY ./src/server/go.mod ./src/server/go.sum ./
RUN go mod download

COPY ./src/server .
RUN mkdir -p ./build
RUN go build -o ./build/server -v -ldflags="-s -w" ./cmd/main.go

COPY ./$STAGE_MODE.env ./build/.env
RUN mkdir -p ./build/misc && cp ./misc/sitemap.json ./build/misc/sitemap.json

FROM node:16.13-alpine as frontbuilder

ARG STAGE_MODE=prod

RUN apk update

WORKDIR /app
RUN mkdir -p ./build_$STAGE_MODE
COPY ./src/server/views ./src/server/views

WORKDIR /app/src/frontend
COPY ./src/frontend/package.json ./src/frontend/package-lock.json ./
RUN npm install

ARG NODE_ENV=production

COPY ./src/frontend .
COPY ./src/server/views ./build_$STAGE_MODE/views
RUN NODE_ENV=$NODE_ENV npm run $STAGE_MODE

FROM alpine:3

ARG STAGE_MODE=prod
ENV TZ=Europe/Moscow

RUN apk update
RUN apk --no-cache add ca-certificates tzdata
RUN update-ca-certificates

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=backbuilder ./app/build /home/app
COPY --from=frontbuilder ./app/build_$STAGE_MODE /home/app
COPY ./data/certbot /etc/letsencrypt

RUN mkdir -p ./home/app/logs
