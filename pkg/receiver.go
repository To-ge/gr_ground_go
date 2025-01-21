package pkg

import (
	"bufio"
	"log"
	"strings"

	"github.com/To-ge/gr_ground_go/config"
	"go.bug.st/serial"
)

type Receiver struct {
	port   serial.Port
	reader *bufio.Reader
}

func NewReceiver() *Receiver {
	mode := &serial.Mode{
		BaudRate: 19200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(config.LoadConfig().Im920sl.UsbFile, mode)
	if err != nil {
		log.Fatalf("serial error: %v\n", err)
	}
	reader := bufio.NewReader(port)

	log.Println("serial port is opened.")
	return &Receiver{
		port:   port,
		reader: reader,
	}
}

func (r *Receiver) Read() (string, error) {
	line, err := r.reader.ReadString('\n') // データを1行ずつ読み取る
	if err != nil {
		log.Printf("Error reading from serial port: %v", err)
		return "", err
	}
	line = strings.TrimSpace(line) // 空白を削除
	return line, nil
}

func (r *Receiver) Close() {
	r.port.Close()
}
