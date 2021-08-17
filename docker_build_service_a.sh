CGO_ENABLED=0 GOOS=linux go build -o service_a KiteX-A/main.go
docker build -t service_a .