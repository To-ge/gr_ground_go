package service

import (
	"context"
	"log"
	"math"
	"time"

	v1 "github.com/To-ge/gr_ground_go/api/gen/go/v1"
	"github.com/To-ge/gr_ground_go/pkg"
)

func LoopSendLocation() {
	conn := pkg.NewGrpcConnection()
	client := v1.NewTelemetryServiceClient(conn.Conn)

	stream, err := client.SendLocation(context.Background())
	if err != nil {
		log.Fatalf("SendLocation stream error: %v\n", err)
	}
	log.Println("SendLocation stream started!")

	for {
		location, ok := <-LocationCh
		if !ok {
			log.Fatalln("LocationCh is closed.")
		}

		timestamp := float64(time.Now().UnixMicro()) / math.Pow10(6)
		pkg.OutputLocationLogger.Printf(",%f,%v,%v,%v\n", timestamp, location.GetLatitude(), location.GetLongitude(), location.GetAltitude())
		if err := stream.Send(location); err != nil {
			log.Printf("SendLocation stream.Send error: %v\n", err)
		}
	}
}
