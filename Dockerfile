FROM golang:1.21 AS builder
RUN ldd --version
WORKDIR /build
COPY . .
RUN go mod download && go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-binary

FROM alpine:3.19
RUN ldd; exit 0
WORKDIR /
COPY --from=builder /build/go-binary .
ENTRYPOINT ["/go-binary"]
