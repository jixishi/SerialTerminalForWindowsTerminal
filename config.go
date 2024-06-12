package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Config struct {
	portName    string
	baudRate    int
	dataBits    int
	stopBits    int
	parityBit   int
	outputCode  string
	inputCode   string
	endStr      string
	enableLog   bool
	logFilePath string
	forWard     []int
	frameSize   int
	timesTamp   bool
	timesFmt    string
	address     []string
}
type FoeWardMode int

const (
	NOT FoeWardMode = iota
	TCPC
	UDPC
)

var config Config

func setForWardClient(mode FoeWardMode, add string) (conn net.Conn) {
	var err error
	switch mode {
	case NOT:

	case TCPC:
		conn, err = net.Dial("tcp", add)
		if err != nil {
			log.Fatal(err)
		}
	case UDPC:
		conn, err = net.Dial("udp", add)
		if err != nil {
			log.Fatal(err)
		}
	default:
		panic("未知模式设置")
	}
	return conn
}

func checkLogOpen() {
	if config.enableLog {
		path := fmt.Sprintf(config.logFilePath, config.portName, time.Now().Format("2006_01_02T150405"))
		f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		outs = append(outs, f)
	}
}
