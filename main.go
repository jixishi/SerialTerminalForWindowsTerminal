package main

import (
	"bufio"
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
	logPath  string
	port     string
	endstr   string
	codein   string
	codeout  string
	baud     int
	prity    int
	stopbits int
	databits int
	logFlag  bool
}
type Cmd struct {
	name string
	des  string
	call func()
}

var (
	conf Config
	cmds []Cmd
	port serial.Port
	err  error

	args []string
)

func CheckPort() {
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
func usage() {
	CheckPort()
	fmt.Printf("\n参数帮助:\n")

	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "p", "port", "", "连接的串口(/dev/ttyUSB0、COMx)", "")
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "b", "baud", 115200, "波特率", 115200)
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "s", "stop", 1, "停止位", 1)
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "d", "data", 8, "数据位", 8)
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "o", "out", "UTF-8", "输出编码", "UTF-8")
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "i", "in", "UTF-8", "输入编码", "UTF-8")
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "e", "end", "\n", "终端换行符", "\\n")
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "l", "log", false, "是否启用日志保存", false)
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "P", "Path", "./Log.txt", "日志保存路径", "./Log.txt")
	fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", "v", "verify", 0, "奇偶校验(0:无校验、1:奇校验、2:偶校验、3:1校验、4:0校验)", 0)
}
func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)
	flag.BoolVar(&conf.logFlag, "log", false, "是否启用日志保存")
	flag.BoolVar(&conf.logFlag, "l", false, "")

	flag.StringVar(&conf.logPath, "Path", "./Log.txt", "日志保存路径")
	flag.StringVar(&conf.logPath, "P", "./Log.txt", "")

	flag.StringVar(&conf.port, "port", "", "连接的串口\t(/dev/ttyUSB0、COMx)")
	flag.StringVar(&conf.port, "p", "", "")

	flag.StringVar(&conf.endstr, "end", "\n", "终端换行符")
	flag.StringVar(&conf.endstr, "e", "\n", "")

	flag.IntVar(&conf.baud, "baud", 115200, "波特率")
	flag.IntVar(&conf.baud, "b", 115200, "")

	flag.IntVar(&conf.prity, "verify", 0, "奇偶校验(0:无校验、1:奇校验、2:偶校验、3:1校验、4:0校验)")
	flag.IntVar(&conf.prity, "v", 0, "")

	flag.IntVar(&conf.stopbits, "stop", 1, "停止位")
	flag.IntVar(&conf.stopbits, "s", 1, "")

	flag.IntVar(&conf.databits, "data", 8, "数据位")
	flag.IntVar(&conf.databits, "d", 8, "")

	flag.StringVar(&conf.codeout, "out", "UTF-8", "输出编码")
	flag.StringVar(&conf.codeout, "o", "UTF-8", "")

	flag.StringVar(&conf.codein, "in", "UTF-8", "输入编码")
	flag.StringVar(&conf.codein, "i", "UTF-8", "")
	cmdinit()
}
func cmdhelp() {
	var page = 0
	fmt.Printf(">-------Help(%v)-------<\n", page)
	for i := 0; i < len(cmds); i++ {
		output(conf.codeout, fmt.Sprintf(" %-10v --%v\n", cmds[i].name, cmds[i].des))
	}
}
func cmdexit() {
	os.Exit(0)
}
func cmdargs() {
	fmt.Printf(">-------Args()-------<\n")
	fmt.Printf("%q\n", args)
}
func cmdinit() {
	cmds = append(cmds, Cmd{name: ".help", des: "帮助信息", call: cmdhelp})
	cmds = append(cmds, Cmd{name: ".args", des: "参数信息", call: cmdargs})
	cmds = append(cmds, Cmd{name: ".exit", des: "退出终端", call: cmdexit})
}
func input() {
	input := bufio.NewScanner(os.Stdin)
	var ok = false
	for {
		input.Scan()
		args = strings.Split(input.Text(), " ")
		for _, cmd := range cmds {
			if strings.Compare(strings.TrimSpace(args[0]), cmd.name) == 0 {
				cmd.call()
				ok = true
			}
		}
		if !ok {
			_, err := port.Write(input.Bytes())
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.WriteString(port, conf.endstr)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
func output(cs, str string) {
	charsetconv.EncodeWith(strings.NewReader(str), os.Stdout, charsetconv.Charset(cs), false)
}
func main() {
	flag.Parse()
	if conf.port == "" {
		fmt.Println("端口未指定")
		usage()
		os.Exit(0)
	}
	mode := &serial.Mode{
		BaudRate: conf.baud,
		StopBits: serial.StopBits(conf.stopbits),
		DataBits: conf.databits,
		Parity:   serial.Parity(conf.prity),
	}
	port, err = serial.Open(conf.port, mode)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()
	go input()
	if conf.logFlag {
		f, err := os.OpenFile(conf.logPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		out := io.MultiWriter(os.Stdout, f)
		for {
			err = charsetconv.ConvertWith(port, charsetconv.Charset(conf.codein), out, charsetconv.Charset(conf.codeout), false)
			//_, err = io.Copy(out, port)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		for {
			err = charsetconv.ConvertWith(port, charsetconv.Charset(conf.codein), os.Stdout, charsetconv.Charset(conf.codeout), false)
			//_, err := io.Copy(os.Stdout, port)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
