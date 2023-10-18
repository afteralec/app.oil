ARG APP_NAME=app

FROM golang:1.19 as builder
ARG APP_NAME
ENV APP_NAME=$APP_NAME
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /$APP_NAME

FROM alpine:latest
ARG APP_NAME
ENV APP_NAME=$APP_NAME
WORKDIR /usr/local/bin/
COPY --from=build /$APP_NAME ./
CMD ./$APP_NAME
