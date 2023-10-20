ARG APP_NAME=oil-app

FROM golang:1.21 as builder
ARG APP_NAME
ENV APP_NAME=$APP_NAME
WORKDIR /app
COPY . .
RUN go mod download
RUN go build cmd/main/main.go -o /$APP_NAME

FROM alpine:latest
ARG APP_NAME
ENV APP_NAME=$APP_NAME
WORKDIR /usr/local/bin/
COPY --from=builder /$APP_NAME ./
CMD ./$APP_NAME
