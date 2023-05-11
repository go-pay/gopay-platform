FROM registry.cn-hangzhou.aliyuncs.com/fmm-ink/golang:1.19.9 AS builder
WORKDIR /build
COPY . /build

RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && go env -w GOSUMDB=off && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w -extldflags '-static'" -o main main.go

FROM registry.cn-hangzhou.aliyuncs.com/fmm-ink/alpine:3.18
WORKDIR /app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && apk add --no-cache curl net-tools busybox-extras iproute2
COPY --from=builder /build/main  /usr/bin/main
COPY --from=builder /build/app/cfg/config.yaml  /app/cfg/config.yaml

CMD ["main", "-conf", "/app/cfg/config.yaml"]