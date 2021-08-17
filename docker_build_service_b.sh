CGO_ENABLED=0 GOOS=linux go build -o service_b KiteX-B/main.go
docker build -t service_b .