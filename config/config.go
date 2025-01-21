package config

import (
	"os"
	"strconv"
)

type appConfig struct {
	Im920sl  im920sl
	GrpcInfo grpcInfo
	TestInfo testInfo
}

type im920sl struct {
	UsbFile string
}

type grpcInfo struct {
	Address         string
	ConnectInterval int
}

type testInfo struct {
	Location location
}

type location struct {
	Latitude  float64
	Longitude float64
	Altitude  float32
}

func LoadConfig() appConfig {
	im920sl := im920sl{
		UsbFile: os.Getenv("USB_DEVICE"),
	}

	interval, _ := strconv.Atoi("GRPC_CONNECT_INTERVAL")
	grpcInfo := grpcInfo{
		Address:         os.Getenv("GRPC_ADDRESS"),
		ConnectInterval: interval,
	}

	latitude, _ := strconv.ParseFloat(os.Getenv("TEST_LATITUDE"), 64)
	longitude, _ := strconv.ParseFloat(os.Getenv("TEST_LONGITUDE"), 64)
	altitude, _ := strconv.ParseFloat(os.Getenv("TEST_ALTITUDE"), 32)
	testInfo := testInfo{
		Location: location{
			Latitude:  latitude,
			Longitude: longitude,
			Altitude:  float32(altitude),
		},
	}

	appConfig := appConfig{
		Im920sl:  im920sl,
		GrpcInfo: grpcInfo,
		TestInfo: testInfo,
	}

	return appConfig
}
