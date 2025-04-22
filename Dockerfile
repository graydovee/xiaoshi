FROM golang:1.22-bullseye as builder

WORKDIR /source/xiaoshi
ENV CGO_ENABLED=1

RUN apt-get update && apt-get install -y gcc sqlite3 libsqlite3-dev

RUN git clone https://github.com/graydovee/ZeroBot-Plugin.git /source/ZeroBot-Plugin

COPY go.mod go.mod
COPY go.sum go.sum

RUN go env -w GOPROXY='https://goproxy.cn,direct'
RUN go mod download


COPY cmd/ cmd/
COPY pkg/ pkg/
COPY main.go main.go

RUN go mod tidy

RUN go build -o /bin/xiaoshi main.go

FROM debian:bullseye-slim
WORKDIR /usr/qqbot

RUN apt-get update && \
    apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /bin/xiaoshi .
ENTRYPOINT ["./xiaoshi"]
