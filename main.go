package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/zimolab/charsetconv"
	"go.bug.st/serial"
	"io"
	"log"
	"os"
	"strings"
)

type Config struct {
	enableLog   bool
	logFilePath string
	portName    string
	endStr      string
	inputCode   string
	outputCode  string
	baudRate    int
	parityBit   int
	stopBits    int
	dataBits    int
}
type Command struct {
	name        string
	description string
	function    func()
}

var (
	config     Config
	commands   []Command
	serialPort serial.Port
	err        error
	args       []string
)

func checkPortAvailability() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("未找到串口！")
	}
	fmt.Printf("找到的串口: ")
	for _, port := range ports {
		fmt.Printf(" %v", port)
	}
}
func printUsage() {
	checkPortAvailability()
	fmt.Printf("\n参数帮助:\n")

	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "p", "portName", "", "连接的串口(/dev/ttyUSB0、COMx)", "")
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "b", "baudRate", 115200, "波特率", 115200)
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "d", "data", 8, "数据位", 8)
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "s", "stop", 0, "停止位停止位(0: 1停止 1:1.5停止 2:2停止)", 0)
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "o", "out", "UTF-8", "输出编码", "UTF-8")
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "i", "in", "UTF-8", "输入编码", "UTF-8")
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "e", "end", "\n", "终端换行符", "\\n")
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "l", "log", false, "是否启用日志保存", false)
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "P", "Path", "./Log.txt", "日志保存路径", "./Log.txt")
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "v", "verify", 0, "奇偶校验(0:无校验、1:奇校验、2:偶校验、3:1校验、4:0校验)", 0)
}
func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)
	flag.BoolVar(&config.enableLog, "log", false, "是否启用日志保存")
	flag.BoolVar(&config.enableLog, "l", false, "")

	flag.StringVar(&config.logFilePath, "Path", "./Log.txt", "日志保存路径")
	flag.StringVar(&config.logFilePath, "P", "./Log.txt", "")

	flag.StringVar(&config.portName, "portName", "", "连接的串口\t(/dev/ttyUSB0、COMx)")
	flag.StringVar(&config.portName, "p", "", "")

	flag.StringVar(&config.endStr, "end", "\n", "终端换行符")
	flag.StringVar(&config.endStr, "e", "\n", "")

	flag.IntVar(&config.baudRate, "baudRate", 115200, "波特率")
	flag.IntVar(&config.baudRate, "b", 115200, "")

	flag.IntVar(&config.parityBit, "verify", 0, "奇偶校验(0:无校验 1:奇校验 2:偶校验 3:1校验 4:0校验)")
	flag.IntVar(&config.parityBit, "v", 0, "")

	flag.IntVar(&config.stopBits, "stop", 0, "停止位(0: 1停止 1:1.5停止 2:2停止)")
	flag.IntVar(&config.stopBits, "s", 0, "")

	flag.IntVar(&config.dataBits, "data", 8, "数据位")
	flag.IntVar(&config.dataBits, "d", 8, "")

	flag.StringVar(&config.outputCode, "out", "UTF-8", "输出编码")
	flag.StringVar(&config.outputCode, "o", "UTF-8", "")

	flag.StringVar(&config.inputCode, "in", "UTF-8", "输入编码")
	flag.StringVar(&config.inputCode, "i", "UTF-8", "")
	cmdinit()
}
func cmdhelp() {
	var page = 0
	fmt.Printf(">-------Help(%v)-------<\n", page)
	for i := 0; i < len(commands); i++ {
		strout(config.outputCode, fmt.Sprintf(" %-10v --%v\n", commands[i].name, commands[i].description))
	}
}
func cmdexit() {
	os.Exit(0)
}
func cmdargs() {
	fmt.Printf(">-------Args(%v)-------<\n", len(args)-1)
	fmt.Printf("%q\n", args[1:])
}
func cmdhex() {
	fmt.Printf(">-----Hex Send-----<\n")
	fmt.Printf("%q\n", args[1:])
	s := strings.Join(args[1:], "")
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}
	_, err = serialPort.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}
func cmdinit() {
	commands = append(commands, Command{name: ".help", description: "帮助信息", function: cmdhelp})
	commands = append(commands, Command{name: ".args", description: "参数信息", function: cmdargs})
	commands = append(commands, Command{name: ".hex", description: "发送Hex", function: cmdhex})
	commands = append(commands, Command{name: ".exit", description: "退出终端", function: cmdexit})
}
func input() {
	input := bufio.NewScanner(os.Stdin)
	var ok = false
	for {
		input.Scan()
		args = strings.Split(input.Text(), " ")
		for _, cmd := range commands {
			if strings.Compare(strings.TrimSpace(args[0]), cmd.name) == 0 {
				cmd.function()
				ok = true
			}
		}
		if !ok {
			_, err := serialPort.Write(input.Bytes())
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.WriteString(serialPort, config.endStr)
			//_, err = io.Copy(portName, os.Stdin)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
func strout(cs, str string) {
	err = charsetconv.EncodeWith(strings.NewReader(str), os.Stdout, charsetconv.Charset(cs), false)
	if err != nil {
		log.Fatal(err)
	}
}
func output(out io.Writer) {
	if strings.Compare(config.inputCode, "hex") == 0 {
		_, err = io.Copy(hex.NewEncoder(out), serialPort)
	} else {
		err = charsetconv.ConvertWith(serialPort, charsetconv.Charset(config.inputCode), out, charsetconv.Charset(config.outputCode), false)
	}
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	flag.Parse()
	if config.portName == "" {
		fmt.Println("端口未指定")
		printUsage()
		os.Exit(0)
	}
	mode := &serial.Mode{
		BaudRate: config.baudRate,
		StopBits: serial.StopBits(config.stopBits),
		DataBits: config.dataBits,
		Parity:   serial.Parity(config.parityBit),
	}
	serialPort, err = serial.Open(config.portName, mode)
	if err != nil {
		log.Fatal(err)
	}
	defer func(port serial.Port) {
		err := port.Close()
		if err != nil {

		}
	}(serialPort)
	go input()
	if config.enableLog {
		f, err := os.OpenFile(config.logFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		out := io.MultiWriter(os.Stdout, f)
		for {
			output(out)
		}
	} else {
		for {
			output(os.Stdout)
		}
	}
}
