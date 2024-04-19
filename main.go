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

var (
	config     Config
	commands   []Command
	serialPort serial.Port
	err        error
	args       []string
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
	cmdinit()
}

func input() {
	input := bufio.NewScanner(os.Stdin)
	var ok = false
	for input.Scan() {
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
		err := serialPort.Drain()
		if err != nil {
			log.Fatal(err)
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
		b, _ := bufio.NewReader(io.LimitReader(serialPort, 16)).Peek(16)
		_, err = fmt.Fprintf(out, "% X %q \n", b, b)
	} else {
		err = charsetconv.ConvertWith(io.LimitReader(serialPort, 1024), charsetconv.Charset(config.inputCode), out, charsetconv.Charset(config.outputCode), false)
	}
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	flag.Parse()
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
