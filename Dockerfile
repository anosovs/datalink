FROM golang:1.20 AS builder
WORKDIR /app
COPY . /app
RUN GO111MODULE=auto CGO_ENABLED=1 GOOS=linux GOPROXY=https://proxy.golang.org go build -ldflags '-w' -o app cmd/main.go

FROM ubuntu:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/templates/ /app/templates/
COPY --from=builder /app/.env /app/.env

ENTRYPOINT ["./app"]
EXPOSE 8080
