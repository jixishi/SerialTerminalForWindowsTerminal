package main

import (
	"log"
	"net"
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
	address     string
}
type FoeWardMode int

const (
	NOT FoeWardMode = iota
	TCPS
	TCPC
	UDPS
	UDPC
)

var config Config

func setForWard() (conn net.Conn) {
	switch FoeWardMode(config.forWard) {
	case TCPS:

	case TCPC:
		conn, err = net.Dial("tcp", config.address)
		if err != nil {
			log.Fatal(err)
		}
	case UDPS:

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
