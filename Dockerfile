ARG APP_NAME=petrichor-app

FROM golang:1.21 as builder
ARG APP_NAME
ENV APP_NAME=$APP_NAME
WORKDIR /
COPY . .
RUN go mod download
RUN go build -o $APP_NAME main.go 

FROM alpine:latest
ARG APP_NAME
ENV APP_NAME=$APP_NAME
WORKDIR /usr/local/bin/
COPY --from=builder /$APP_NAME ./
CMD ./$APP_NAME
