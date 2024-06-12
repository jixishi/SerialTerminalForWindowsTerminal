package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)
	for _, f := range flags {
		flagInit(&f)
	}
	cmdinit()
}

func main() {
	pflag.Parse()
	flagExt()
	if config.portName == "" {
		getCliFlag()
	}
	ports, err := checkPortAvailability(config.portName)
	if err != nil {
		fmt.Println(err)
		printUsage(ports)
		os.Exit(0)
	}

	// 日志文件输出检测
	checkLogOpen()

	//串口设备开启
	OpenSerial()

	defer CloseSerial()
	// 打开文件服务
	OpenTrzsz()

	defer CloseTrzsz()

	//开启转发
	OpenForwarding()

	// 获取终端输入
	go input(in)

	if len(outs) != 1 {
		out = io.MultiWriter(outs...)
	}
	for {
		output()
	}
}
