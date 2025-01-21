package pkg

import (
	"log"
	"time"

	"github.com/To-ge/gr_ground_go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GrpcConnection struct {
	Conn *grpc.ClientConn
}

func NewGrpcConnection() GrpcConnection {
	address := config.LoadConfig().GrpcInfo.Address
	interval := config.LoadConfig().GrpcInfo.ConnectInterval

	for {
		creds := credentials.NewTLS(nil)
		conn, err := grpc.NewClient(
			address,
			grpc.WithTransportCredentials(creds),
		)
		if err != nil {
			log.Fatal("gRPC Connection failed.")
			time.Sleep(time.Duration(interval) * time.Second)
			continue
		}
		log.Println("gRPC Connection started.")

		return GrpcConnection{
			Conn: conn,
		}
	}
}

func (c GrpcConnection) Close() {
	c.Conn.Close()
}
