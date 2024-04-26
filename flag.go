package main

import (
	"flag"
	"fmt"
	"strings"
)

type ptrVal struct {
	*string
	*int
	*bool
	*float64
	*float32
}
type Val struct {
	string
	int
	bool
	float64
	float32
}
type Flag struct {
	v    ptrVal
	sStr string
	lStr string
	dv   Val
	help string
}

var (
	portName    = Flag{ptrVal{string: &config.portName}, "p", "port", Val{string: ""}, "要连接的串口\t(/dev/ttyUSB0、COMx)"}
	baudRate    = Flag{ptrVal{int: &config.baudRate}, "b", "baud", Val{int: 115200}, "波特率"}
	dataBits    = Flag{ptrVal{int: &config.dataBits}, "d", "data", Val{int: 8}, "数据位"}
	stopBits    = Flag{ptrVal{int: &config.stopBits}, "s", "stop", Val{int: 0}, "停止位停止位(0: 1停止 1:1.5停止 2:2停止)"}
	outputCode  = Flag{ptrVal{string: &config.outputCode}, "o", "out", Val{string: "UTF-8"}, "输出编码"}
	inputCode   = Flag{ptrVal{string: &config.inputCode}, "i", "in", Val{string: "UTF-8"}, "输入编码"}
	endStr      = Flag{ptrVal{string: &config.endStr}, "e", "end", Val{string: "\n"}, "终端换行符"}
	enableLog   = Flag{ptrVal{bool: &config.enableLog}, "l", "log", Val{bool: false}, "是否启用日志保存"}
	logFilePath = Flag{ptrVal{string: &config.logFilePath}, "P", "Path", Val{string: "./Log.txt"}, "日志保存路径"}
	forWard     = Flag{ptrVal{int: &config.forWard}, "f", "forward", Val{int: 0}, "转发模式(0: 无 1:TCP-C 2:UDP-C)"}
	address     = Flag{ptrVal{string: &config.address}, "a", "address", Val{string: "127.0.0.1:12345"}, "转发服务地址"}
	frameSize   = Flag{ptrVal{int: &config.forWard}, "F", "Frame", Val{int: 16}, "帧大小"}
	parityBit   = Flag{ptrVal{int: &config.parityBit}, "v", "verify", Val{int: 0}, "奇偶校验(0:无校验、1:奇校验、2:偶校验、3:1校验、4:0校验)"}
	flags       = []Flag{portName, baudRate, dataBits, stopBits, outputCode, inputCode, endStr, enableLog, logFilePath, forWard, frameSize, address, parityBit}
)

type ValType int

const (
	notVal ValType = iota
	stringVal
	intVal
	boolVal
)

func printUsage(ports []string) {
	fmt.Printf("\n参数帮助:\n")
	for _, f := range flags {
		flagprint(f)
	}
	fmt.Printf("\n在线串口: %v\n", strings.Join(ports, ","))
}
func flagFindValue(v ptrVal) ValType {
	if v.string != nil {
		return stringVal
	}
	if v.bool != nil {
		return boolVal
	}
	if v.int != nil {
		return intVal
	}
	return notVal
}
func flagprint(f Flag) {
	switch flagFindValue(f.v) {
	case stringVal:
		fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%q\n", f.sStr, f.lStr, f.dv.string, f.help, f.dv.string)
	case intVal:
		fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", f.sStr, f.lStr, f.dv.int, f.help, f.dv.int)
	case boolVal:
		fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", f.sStr, f.lStr, f.dv.bool, f.help, f.dv.bool)
	default:
		panic("unhandled default case")
	}
}
func flagInit(f *Flag) {
	if f.v.string != nil {
		flag.StringVar(f.v.string, f.sStr, f.dv.string, "")
		flag.StringVar(f.v.string, f.lStr, f.dv.string, f.help)
	}
	if f.v.bool != nil {
		flag.BoolVar(f.v.bool, f.sStr, f.dv.bool, "")
		flag.BoolVar(f.v.bool, f.lStr, f.dv.bool, f.help)
	}
	if f.v.int != nil {
		flag.IntVar(f.v.int, f.sStr, f.dv.int, "")
		flag.IntVar(f.v.int, f.lStr, f.dv.int, f.help)
	}
}
