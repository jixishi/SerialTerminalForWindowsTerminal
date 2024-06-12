package main

import (
	"fmt"
	"github.com/trzsz/trzsz-go/trzsz"
	"go.bug.st/serial"
	"golang.org/x/term"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
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

func OpenSerial() {
	var err error
	mode := &serial.Mode{
		BaudRate: config.baudRate,
		StopBits: serial.StopBits(config.stopBits),
		DataBits: config.dataBits,
		Parity:   serial.Parity(config.parityBit),
	}
	serialPort, err = serial.Open(config.portName, mode)
	ErrorF(err)
	return
}

func CloseSerial() {
	err := serialPort.Close()
	ErrorF(err)
	return
}

var termch chan os.Signal

// OpenTrzsz create a TrzszFilter to support trzsz ( trz / tsz ).
//
// ┌────────┐   stdinPipe  ┌────────┐   ClientIn   ┌─────────────┐   SerialIn   ┌────────┐
// │        ├─────────────►│        ├─────────────►│             ├─────────────►│        │
// │ mutual │              │ Client │              │ TrzszFilter │              │ Serial │
// │        │◄─────────────│        │◄─────────────┤             │◄─────────────┤        │
// └────────┘   stdoutPipe └────────┘   ClientOut  └─────────────┘   SerialOut  └────────┘
func OpenTrzsz() {
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
	termch = make(chan os.Signal, 1)
	go func() {
		for range termch {
			width, _, err := term.GetSize(fd)
			if err != nil {
				fmt.Printf("term get size failed: %s\n", err)
				continue
			}
			trzszFilter.SetTerminalColumns(int32(width))
		}
	}()
}

func CloseTrzsz() {
	signal.Stop(termch)
	close(termch)
}

func OpenForwarding() {
	for i, mode := range config.forWard {
		if FoeWardMode(mode) != NOT {
			conn := setForWardClient(FoeWardMode(mode), config.address[i])
			outs = append(outs, conn)
			go func() {
				defer func(conn net.Conn) {
					err := conn.Close()
					if err != nil {
						log.Fatal(err)
					}
				}(conn)
				input(conn)
			}()
		}
	}
}

func ErrorP(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}
func ErrorF(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
