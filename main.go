package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/zimolab/charsetconv"
	"go.bug.st/serial"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var (
	serialPort serial.Port
	err        error
	args       []string
)

var (
	in   io.Reader = os.Stdin
	out  io.Writer = os.Stdout
	ins            = []io.Reader{os.Stdin}
	outs           = []io.Writer{os.Stdout}
)

func checkPortAvailability(name string) ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		return nil, fmt.Errorf("无串口")
	}
	if name == "" {
		return ports, fmt.Errorf("串口未指定")
	}
	for _, port := range ports {
		if strings.Compare(port, name) == 0 {
			return ports, nil
		}
	}
	return ports, fmt.Errorf("串口 " + name + " 未在线")
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)
	for _, f := range flags {
		flagInit(&f)
	}
	flag.Func("h", "获取帮助", func(s string) error {
		ports, err := checkPortAvailability(s)
		if err != nil {
			fmt.Println(err)
			printUsage(ports)
			os.Exit(0)
		}
		return err
	})
	cmdinit()
}

func input(in io.Reader) {
	input := bufio.NewScanner(in)
	var ok = false
	for {
		input.Scan()
		ok = false
		args = strings.Split(input.Text(), " ")
		for _, cmd := range commands {
			if strings.Compare(strings.TrimSpace(args[0]), cmd.name) == 0 {
				cmd.function()
				ok = true
			}
		}
		if !ok {
			_, err := io.WriteString(serialPort, input.Text())
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.WriteString(serialPort, config.endStr)
			if err != nil {
				log.Fatal(err)
			}
		}
		err = serialPort.Drain()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func strout(out io.Writer, cs, str string) {
	err = charsetconv.EncodeWith(strings.NewReader(str), out, charsetconv.Charset(cs), false)
	if err != nil {
		log.Fatal(err)
	}
}

func output() {
	if strings.Compare(config.inputCode, "hex") == 0 {
		b := make([]byte, 16)
		r, _ := io.LimitReader(serialPort, int64(config.frameSize)).Read(b)
		if r != 0 {
			strout(out, config.outputCode, fmt.Sprintf("% X %q \n", b, b))
		}
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
		getCliFlag()
	}
	ports, err := checkPortAvailability(config.portName)
	if err != nil {
		fmt.Println(err)
		printUsage(ports)
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
			log.Fatal(err)
		}
	}(serialPort)

	if FoeWardMode(config.forWard) != NOT {
		conn := setForWardClient()
		ins = append(ins, conn)
		outs = append(outs, conn)
		defer func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(conn)
	}
	if len(ins) != 0 {
		for _, reader := range ins {
			go input(reader)
		}
	}
	if config.enableLog {
		f, err := os.OpenFile(config.logFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		outs = append(outs, f)
	}
	if len(outs) != 1 {
		out = io.MultiWriter(outs...)
	}
	for {
		output()
	}
}
