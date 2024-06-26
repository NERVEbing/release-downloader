FROM golang:1.21.10 AS builder

COPY . /app

WORKDIR /app

ENV GO111MODULE=on \
    CGO_ENABLED=0

RUN mkdir -p ./bin && go build -o ./bin ./...

FROM alpine:3.19.1

RUN apk update && apk --no-cache add tzdata

ENV TZ=Asia/Shanghai

COPY --from=builder /app/bin /app

WORKDIR /app

CMD ["./release-downloader"]
