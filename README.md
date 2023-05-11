# gopay-platform

GoPay payment platform

## Terminal Run

```bash
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w -extldflags '-static'" -o main main.go

main -conf app/cfg/config.yaml
```