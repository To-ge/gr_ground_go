package pkg

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

var (
	InputLocationLogger  *log.Logger
	OutputLocationLogger *log.Logger
	Bme280Logger         *log.Logger
	Mpu6050AccelLogger   *log.Logger
	Mpu6050GyroLogger    *log.Logger
	Im920slLogger        *log.Logger
)

func InitLogger() {
	if err := createLogFolder("../log"); err != nil {
		log.Fatal(err.Error())
	}
	rl, err := rotatelogs.New(
		"../log/server_%Y%m%d.log",
		rotatelogs.WithRotationTime(time.Hour*time.Duration(24)), //rotation time
		rotatelogs.WithRotationCount(7),                          //max backup count
	)
	fmt.Println(os.Getwd())
	if err != nil {
		fmt.Println("error: rotatelog.New()")
		return
	}
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(rl)
}

func InitTimestampLogger() {
	today := time.Now().Format("20060102")

	// GPS
	if err := createLogFolder("../log/gps"); err != nil {
		log.Fatal(err.Error())
	}
	inputLocationLogFile, err := os.OpenFile(fmt.Sprintf("../log/gps/gps-in_%s.log", today), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open ../log/gps/gps-in_%s.log file: %s", today, err.Error())
	}
	outputLocationLogFile, err := os.OpenFile(fmt.Sprintf("../log/gps/gps-out_%s.log", today), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open ../log/gps/gps-out_%s.log file: %s", today, err.Error())
	}

	// BME280
	if err := createLogFolder("../log/bme280"); err != nil {
		log.Fatal(err.Error())
	}
	bme280LogFile, err := os.OpenFile(fmt.Sprintf("../log/bme280/bme280_%s.log", today), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open ../log/bme280/bme280_%s.log file: %s", today, err.Error())
	}

	// MPU6050
	if err := createLogFolder("../log/mpu6050"); err != nil {
		log.Fatal(err.Error())
	}
	mpu6050AccelLogFile, err := os.OpenFile(fmt.Sprintf("../log/mpu6050/mpu6050-accel_%s.log", today), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open ../log/mpu6050/mpu6050-accel_%s.log file: %s", today, err.Error())
	}
	mpu6050GyroLogFile, err := os.OpenFile(fmt.Sprintf("../log/mpu6050/mpu6050-gyro_%s.log", today), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open ../log/mpu6050/mpu6050-gyro_%s.log file: %s", today, err.Error())
	}

	// IM920sL
	if err := createLogFolder("../log/im920sl"); err != nil {
		log.Fatal(err.Error())
	}
	im920slLogFile, err := os.OpenFile(fmt.Sprintf("../log/im920sl/im920sl_%s.log", today), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open ../log/im920sl/im920sl_%s.log file: %s", today, err.Error())
	}

	// ログ出力を分ける
	InputLocationLogger = log.New(io.MultiWriter(inputLocationLogFile, os.Stdout), "", log.LstdFlags)
	OutputLocationLogger = log.New(io.MultiWriter(outputLocationLogFile, os.Stderr), "", log.LstdFlags)
	Bme280Logger = log.New(io.MultiWriter(bme280LogFile, os.Stderr), "", log.LstdFlags)
	Mpu6050AccelLogger = log.New(io.MultiWriter(mpu6050AccelLogFile, os.Stderr), "", log.LstdFlags)
	Mpu6050GyroLogger = log.New(io.MultiWriter(mpu6050GyroLogFile, os.Stderr), "", log.LstdFlags)
	Im920slLogger = log.New(io.MultiWriter(im920slLogFile, os.Stderr), "", log.LstdFlags)
}

func createLogFolder(folderName string) error {
	_, err := os.Stat(folderName)
	if os.IsNotExist(err) {
		err := os.Mkdir(folderName, os.ModePerm)
		if err != nil {
			return fmt.Errorf("フォルダの作成に失敗しました: %w", err)
		}
		fmt.Println("フォルダを作成しました:", folderName)
	} else if err != nil {
		return fmt.Errorf("フォルダの状態確認に失敗しました: %w", err)
	} else {
		fmt.Println("フォルダは既に存在します:", folderName)
	}
	return nil
}
