FROM golang:latest as builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
EXPOSE 8086
RUN GOOS=linux GOARCH=amd64 go build -o /app/main /app/cmd/app/main.go

FROM ubuntu:18.04

COPY --from=builder /app/main /usr/bin/main
COPY --from=builder /app/logs /usr/bin/logs
COPY --from=builder /app/docker/env /usr/bin/.env
WORKDIR /usr/bin/
EXPOSE 8086
ENTRYPOINT [ "/usr/bin/main" ]