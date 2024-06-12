package main

import (
	"bufio"
	"fmt"
	"github.com/trzsz/trzsz-go/trzsz"
	"github.com/zimolab/charsetconv"
	"go.bug.st/serial"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var (
	serialPort  serial.Port
	in          io.Reader = os.Stdin
	out         io.Writer = os.Stdout
	outs                  = []io.Writer{os.Stdout}
	trzszFilter *trzsz.TrzszFilter
	clientIn    *io.PipeReader
	stdoutPipe  *io.PipeReader
	stdinPipe   *io.PipeWriter
	clientOut   *io.PipeWriter
)

func input(in io.Reader) {
	var err error
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
		ErrorF(err)
	}
}

func strout(out io.Writer, cs, str string) {
	err := charsetconv.EncodeWith(strings.NewReader(str), out, charsetconv.Charset(cs), false)
	ErrorF(err)
}

func output() {
	var err error
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
	ErrorP(err)
}
