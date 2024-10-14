FROM golang:1.22.5 AS builder
WORKDIR /build

# 设置编译环境变量
# GODEBUG=tlsrsakex=1。By default, cipher suites without ECDHE support are no longer offered by either clients or servers during pre-TLS 1.3 handshakes. This change can be reverted with the tlsrsakex=1 GODEBUG setting.
# https://go.dev/doc/go1.22
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct \
    GOPRIVATE=gomod.sunmi.com \
    GOSUMDB=off \
    GODEBUG=tlsrsakex=1

# 将 go.mod 复制到容器中
COPY go.mod ./

# 下载依赖模块，并使用缓存
RUN --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# 把当前目录的所有内容copy到 WORKDIR指定的目录中
COPY . .

# Go编译，并使用缓存
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o main main.go && strip main

FROM alpine:3.19.1
WORKDIR /app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && apk add --no-cache curl net-tools busybox-extras iproute2
COPY --from=builder /build/main  /usr/bin/main
COPY --from=builder /build/app/conf/config.yaml  /app/conf/config.yaml

CMD ["main", "-conf", "/app/conf/config.yaml"]