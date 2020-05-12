package main

import (
        "context"
        "flag"
        "fmt"
        "net"

        log "github.com/sirupsen/logrus"
        "google.golang.org/grpc"
        "google.golang.org/grpc/reflection"
        accesslog "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v2"
)

var (
        localhost = "127.0.0.1"
        alsPort   uint
)

func init() {
        flag.UintVar(&alsPort, "als", 18090, "Accesslog server port")

}

func main() {
        flag.Parse()

        ctx := context.Background()
        log.Printf("Starting grpc access log server")

        als := &AccessLogService{}
        RunAccessLogServer(ctx, als, alsPort)

}

// RunAccessLogServer starts an accesslog service.
func RunAccessLogServer(ctx context.Context, als *AccessLogService, port uint) {
        grpcServer := grpc.NewServer()

        accesslog.RegisterAccessLogServiceServer(grpcServer, als)
        reflection.Register(grpcServer)
      
        lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
        if err != nil {
                log.Errorf("error in listening to port ", err)
        }
        log.Printf("Listening at ", port)
        if err = grpcServer.Serve(lis); err != nil {
                        log.Errorf("server error",err)
                }
      
}
