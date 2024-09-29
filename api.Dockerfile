FROM golang:1.20.5-alpine AS builder

ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main


FROM alpine:3.18.4

RUN apk add --no-cache curl tzdata ca-certificates \
    && rm -rf /var/cache/apk/*

ARG APP_VERSION

ENV APP_VERSION $APP_VERSION
ENV TZ Asia/Jakarta
ENV APP_PORT 80

WORKDIR /app

COPY --from=builder ./app/main .

CMD ["./main"]
