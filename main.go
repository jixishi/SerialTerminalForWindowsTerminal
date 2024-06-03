package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/trzsz/trzsz-go/trzsz"
	"github.com/zimolab/charsetconv"
	"go.bug.st/serial"
	"golang.org/x/term"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"
)

var (
	serialPort serial.Port
	err        error
	args       []string
)

var (
	in          io.Reader = os.Stdin
	out         io.Writer = os.Stdout
	ins                   = []io.Reader{os.Stdin}
	outs                  = []io.Writer{os.Stdout}
	trzszFilter *trzsz.TrzszFilter
	clientIn    *io.PipeReader
	stdoutPipe  *io.PipeReader
	stdinPipe   *io.PipeWriter
	clientOut   *io.PipeWriter
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
			_, err := io.WriteString(stdinPipe, input.Text())
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.WriteString(stdinPipe, config.endStr)
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
		b := make([]byte, config.frameSize)
		r, _ := io.LimitReader(stdoutPipe, int64(config.frameSize)).Read(b)
		if r != 0 {
			if config.timesTamp {
				strout(out, config.outputCode, fmt.Sprintf("%v % X %q \n", time.Now().Format(config.timesFmt), b, b))
			} else {
				strout(out, config.outputCode, fmt.Sprintf("% X %q \n", b, b))
			}
		}
	} else {
		if config.timesTamp {
			line, _, _ := bufio.NewReader(stdoutPipe).ReadLine()
			if line != nil {
				strout(out, config.outputCode, fmt.Sprintf("%v %s\n", time.Now().Format(config.timesFmt), line))
			}
		} else {
			err = charsetconv.ConvertWith(stdoutPipe, charsetconv.Charset(config.inputCode), out, charsetconv.Charset(config.outputCode), false)
		}
	}
	if err != nil {
		log.Fatal(err)
	}
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
	fd := int(os.Stdin.Fd())
	width, _, err := term.GetSize(fd)
	if err != nil {
		if runtime.GOOS != "windows" {
			fmt.Printf("term get size failed: %s\n", err)
			return
		}
		width = 80
	}

	clientIn, stdinPipe = io.Pipe()
	stdoutPipe, clientOut = io.Pipe()
	trzszFilter = trzsz.NewTrzszFilter(clientIn, clientOut, serialPort, serialPort,
		trzsz.TrzszOptions{TerminalColumns: int32(width), EnableZmodem: true})
	trzsz.SetAffectedByWindows(false)
	ch := make(chan os.Signal, 1)
	go func() {
		for range ch {
			width, _, err := term.GetSize(fd)
			if err != nil {
				fmt.Printf("term get size failed: %s\n", err)
				continue
			}
			trzszFilter.SetTerminalColumns(int32(width))
		}
	}()
	defer func() { signal.Stop(ch); close(ch) }()

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
	checkLogOpen()

	if len(ins) != 0 {
		for _, reader := range ins {
			go input(reader)
		}
	}
	if len(outs) != 1 {
		out = io.MultiWriter(outs...)
	}
	for {
		output()
	}
}
