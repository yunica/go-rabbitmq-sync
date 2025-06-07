FROM golang:1.24.4-alpine3.22 AS builder
ENV GO111MODULE=on

RUN mkdir -p /app/bin
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

