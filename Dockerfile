FROM golang:1.21

LABEL maintainer="Nicholas Arvelo nicholas@arvelo.dev"
LABEL com.na_ddns.version="1.0"
LABEL com.na_ddns.description="na.DDNS - Cloudflare Dynamic DNS Client"
LABEL com.na_ddns.release-date="2023-10-01"

WORKDIR /usr/src/na.ddns

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/na.ddns ./...

CMD ["na.ddns"]

