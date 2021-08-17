CGO_ENABLED=0 GOOS=linux go build -o service_c Kitex-C/main.go
docker build -t service_c .