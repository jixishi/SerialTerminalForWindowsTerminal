package main

import (
	"log"
	"net"
	"os"
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
	forWard     int
	frameSize   int
	timesTamp   bool
	timesFmt    string
	address     string
}
type FoeWardMode int

const (
	NOT FoeWardMode = iota
	TCPC
	UDPC
)

var config Config

func setForWardClient() (conn net.Conn) {
	switch FoeWardMode(config.forWard) {
	case NOT:

	case TCPC:
		conn, err = net.Dial("tcp", config.address)
		if err != nil {
			log.Fatal(err)
		}
	case UDPC:
		conn, err = net.Dial("udp", config.address)
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
		f, err := os.OpenFile(config.logFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		outs = append(outs, f)
	}
}
