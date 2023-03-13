FROM golang:1.19.7-alpine3.17 as builder
WORKDIR /app

RUN apk update
RUN apk add libc6-compat git build-base

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd ./cmd
COPY ./cfg ./cfg
COPY ./pkg ./pkg
COPY ./docs ./docs
COPY example.env .

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go build -o main ./cmd/main.go

FROM alpine:latest as runner
WORKDIR /app

RUN apk add libc6-compat

ENV GRPC_PORT=8020
ENV HTTP_PORT=8021

ENV TOKEN_SECRET=32-characters-some-strong-secret

ENV DEFAULT_TOKEN_DURATION=10m
ENV MAXIMUM_TOKEN_DURATION=168h

ENV API_HEADER=api-secret
ENV API_SECRET=tokens

EXPOSE 8020
EXPOSE 8021

COPY --from=builder /app/main .
COPY --from=builder /app/example.env .

ENTRYPOINT [ "/app/main" ]
