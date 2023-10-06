FROM golang:latest
LABEL maintainer="nicholas@arvelo.dev"

ARG BUILD_DATE
ARG BUILD_VERSION

LABEL org.label-schema.schema-version="1.0"

LABEL org.label-schema.build-date="$BUILD_DATE"
LABEL org.label-schema.version="$BUILD_VERSION"

LABEL org.label-schema.name="na.DDNS"
LABEL org.label-schema.description="Cloudflare Dynamic DNS Client"
LABEL org.label-schema.vcs-url="https://github.com/nicholasarvelo/na.DDNS"

WORKDIR /usr/src/na.ddns

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/na.ddns ./...

CMD ["na.ddns"]

