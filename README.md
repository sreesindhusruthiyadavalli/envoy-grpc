# envoy-grpc
#To build the app
CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o main .

#To build docker image 
docker build -t <tag> .
