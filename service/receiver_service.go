package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	v1 "github.com/To-ge/gr_ground_go/api/gen/go/v1"
	"github.com/To-ge/gr_ground_go/pkg"
)

type Telemetry uint

const (
	Unknown Telemetry = iota
	Gps
	Bme
	MpuAccel
	MpuGyro
)

var (
	LocationCh = make(chan *v1.SendLocationRequest)
	testData   = []string{
		"00,0001,D6:1A,0F,0B,0F,0C,0F,00",
		"00,0001,D9:1A,31,F5,72,36,B1,30,F5,42,72,C1,7F,40",
		"00,0001,D7:3A,1F,00,8B,0F,02,6C,0F,09,70",
		"00,0001,D7:4A,E1,F7,63,BE,0F,65,6C,E0,F4,58",
		"00,0001,D7:3A,1F,00,8B,0F,01,9C,0F,09,10",
		"00,0001,D6:4A,E1,F6,41,BE,0F,83,2C,E0,F4,43",
		"00,0001,D6:3A,1F,00,2B,0F,02,4C,0F,08,60",
		"00,0001,D7:4A,E1,F5,27,BE,0F,55,CE,0F,33,60",
	}
)

func LoopReceive() {
	receiver := pkg.NewReceiver()

	// index := 0
	// for len(testData) > index {
	for {
		line, err := receiver.Read()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		// line := testData[index]
		if len(line) == 0 {
			continue
		}
		pkg.Im920slLogger.Println(line)

		_, payload, err := divideData(line)
		if err != nil {
			continue
		}

		formatedPayload := formatPayload(payload)
		telemetry, data := formatTelemetry(formatedPayload)
		if telemetry == Gps {
			LocationCh <- &v1.SendLocationRequest{
				Timestamp: data["timestamp"].(int64),
				Latitude:  data["latitude"].(float64),
				Longitude: data["longitude"].(float64),
				Altitude:  data["altitude"].(float32),
			}
		}

		// index++
	}
}

func divideData(data string) (header string, payload string, err error) {
	twoData := strings.Split(data, ":")
	if len(twoData) != 2 {
		return "", "", fmt.Errorf("data format is invalid")
	}
	return twoData[0], twoData[1], nil
}

func formatPayload(payload string) string {
	payload = strings.Join(strings.Split(payload, ","), "")
	payload = strings.ReplaceAll(payload, "E", "-")
	payload = strings.ReplaceAll(payload, "F", ".")
	return payload
}

func formatTelemetry(data string) (Telemetry, map[string]interface{}) {
	if len(data) < 1 {
		return Unknown, nil
	}

	telemetryType, err := strconv.Atoi(data[0:1])
	if err != nil {
		log.Printf("Invalid telemetry type: %v", err)
		return Unknown, nil
	}
	content := data[1:]

	switch Telemetry(telemetryType) {
	case Gps:
		if !checkSymbols("ABC", content) {
			log.Println("Required symbols not found in GPS data")
			return Unknown, nil
		}
		latitude := strings.Split(strings.Split(content, "A")[1], "B")[0]
		longitude := strings.Split(strings.Split(content, "B")[1], "C")[0]
		altitude := strings.Split(content, "C")[1]

		latitude64, _ := strconv.ParseFloat(latitude, 64)
		longitude64, _ := strconv.ParseFloat(longitude, 64)
		altitude32, _ := (strconv.ParseFloat(altitude, 32))
		timestamp := float64(time.Now().UnixMilli()) / 1000.0

		pkg.InputLocationLogger.Printf(",%f,%v,%v,%v\n", timestamp, latitude64, longitude64, altitude32)

		return Gps, map[string]interface{}{
			"timestamp": int64(timestamp),
			"latitude":  latitude64,
			"longitude": longitude64,
			"altitude":  float32(altitude32),
		}
	case Bme:
		if !checkSymbols("ABCD", content) {
			log.Println("Required symbols not found in BME280 data")
			return Unknown, nil
		}
		temperature := strings.Split(strings.Split(content, "A")[1], "B")[0]
		humidity := strings.Split(strings.Split(content, "B")[1], "C")[0]
		pressure := strings.Split(strings.Split(content, "C")[1], "D")[0]
		altitude := strings.Split(content, "D")[1]

		temperature64, _ := strconv.ParseFloat(temperature, 64)
		humidity64, _ := strconv.ParseFloat(humidity, 64)
		pressure64, _ := strconv.ParseFloat(pressure, 64)
		altitude64, _ := strconv.ParseFloat(altitude, 64)

		pkg.Bme280Logger.Printf(",%v,%v,%v,%v\n", temperature64, humidity64, pressure64, altitude64)

		return Bme, map[string]interface{}{
			"temperature": temperature64,
			"humidity":    humidity64,
			"pressure":    pressure64,
			"altitude":    altitude64,
		}
	case MpuAccel:
		if !checkSymbols("ABC", content) {
			log.Println("Required symbols not found in Mpu6050-Accel data")
			return Unknown, nil
		}
		accelX := strings.Split(strings.Split(content, "A")[1], "B")[0]
		accelY := strings.Split(strings.Split(content, "B")[1], "C")[0]
		accelZ := strings.Split(content, "C")[1]

		accelX64, _ := strconv.ParseFloat(accelX, 64)
		accelY64, _ := strconv.ParseFloat(accelY, 64)
		accelZ64, _ := strconv.ParseFloat(accelZ, 64)

		pkg.Mpu6050AccelLogger.Printf(",%v,%v,%v\n", accelX64, accelY64, accelZ64)

		return MpuAccel, map[string]interface{}{
			"accel_x": accelX64,
			"accel_y": accelY64,
			"accel_z": accelZ64,
		}
	case MpuGyro:
		if !checkSymbols("ABC", content) {
			log.Println("Required symbols not found in Mpu6050-Gyro data")
			return Unknown, nil
		}
		gyroX := strings.Split(strings.Split(content, "A")[1], "B")[0]
		gyroY := strings.Split(strings.Split(content, "B")[1], "C")[0]
		gyroZ := strings.Split(content, "C")[1]

		gyroX64, _ := strconv.ParseFloat(gyroX, 64)
		gyroY64, _ := strconv.ParseFloat(gyroY, 64)
		gyroZ64, _ := strconv.ParseFloat(gyroZ, 64)

		pkg.Mpu6050GyroLogger.Printf(",%v,%v,%v\n", gyroX64, gyroY64, gyroZ64)

		return MpuGyro, map[string]interface{}{
			"gyro_x": gyroX64,
			"gyro_y": gyroY64,
			"gyro_z": gyroZ64,
		}
	default:
		log.Println("Unknown telemetry type")
		return Unknown, nil
	}
}

func checkSymbols(symbols, target string) bool {
	for _, s := range symbols {
		if strings.ContainsRune(target, s) {
			continue
		} else {
			return false
		}
	}
	return true
}
