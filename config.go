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
	frameSize   int
	address     string
}
type FoeWardMode int

const (
	NOT FoeWardMode = iota
	TCPS
	TCPC
	UDPC
)

var config Config

func setForWardClient() (conn net.Conn) {
	switch FoeWardMode(config.forWard) {
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
func setForWardServer() {
	switch FoeWardMode(config.forWard) {
	case TCPS:
		listen, err := net.Listen("tcp", config.address)
		if err != nil {
			log.Fatal(err)
		}
		for {
			conn, err := listen.Accept() // 监听客户端的连接请求
			if err != nil {
				log.Println("Accept() failed, err: ", err)
				continue
			}
			go process(conn) // 启动一个goroutine来处理客户端的连接请求
		}
	default:
		panic("未知模式设置")
	}
}
func process(conn net.Conn) {
	defer conn.Close() // 关闭连接
	//reader := bufio.NewReader(serialPort)
	outs = append(outs, conn)
	defer func() {
		for i, w := range outs {
			if w == conn {
				outs = append(outs[:i], outs[i+1:]...)
			}
		}
	}()
	input(conn)
}
