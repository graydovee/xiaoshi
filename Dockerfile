FROM golang:1.23 as builder

WORKDIR /app
ENV CGO_ENABLED=0

RUN go env -w GOPROXY='https://goproxy.cn,direct'

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN make build

FROM ubuntu:22.04
WORKDIR /app/xiaoshi

RUN apt-get update

# base environment
RUN apt-get install -y curl build-essential

# install python3
RUN apt-get install -y python3

# install nodejs
RUN (curl -fsSL https://deb.nodesource.com/setup_22.x | bash -) && apt-get install -y nodejs

COPY --from=builder /app/bin/xiaoshi .
COPY --from=builder /app/mcp-server ./mcp-server
ENTRYPOINT ["./xiaoshi"]
