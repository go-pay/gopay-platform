# gopay-platform

GoPay payment platform

## Terminal Run

```bash
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go && strip main

main -conf app/conf/config.yaml
```