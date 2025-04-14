FROM golang:1.23.3-alpine AS builder

RUN apk update && apk add ca-certificates git gcc g++ libc-dev binutils

WORKDIR /opt

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o bin/bot ./cmd/bot

FROM alpine:3.19 AS runner

RUN apk update && apk add ca-certificates libc6-compat openssh bash && rm -rf /var/cache/apk/*

WORKDIR /opt

COPY --from=builder /opt/bin/bot ./

COPY telegram.yml /opt/telegram.yml
COPY locales /opt/locales
COPY client_banner.png /opt/client_banner.png
COPY auth_banner.png /opt/auth_banner.png
COPY advertiser_banner.png /opt/advertiser_banner.png

CMD ["./bot"]