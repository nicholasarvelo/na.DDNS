FROM golang:latest

LABEL org.opencontainers.title="na.DDNS"
LABEL org.opencontainers.image.authors="nicholas@arvelo.dev"
LABEL org.opencontainers.description="Cloudflare Dynamic DNS Client"
LABEL org.opencontainers.source="https://github.com/nicholasarvelo/na.DDNS"

WORKDIR /usr/src/na.ddns

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/na.ddns ./...

CMD ["na.ddns"]

