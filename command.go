package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
)

type Command struct {
	name        string
	description string
	function    func()
}

var commands []Command

func cmdhelp() {
	var page = 0
	strout(out, config.outputCode, fmt.Sprintf(">-------Help(%v)-------<\n", page))
	for i := 0; i < len(commands); i++ {
		strout(out, config.outputCode, fmt.Sprintf(" %-10v --%v\n", commands[i].name, commands[i].description))
	}
}
func cmdexit() {
	os.Exit(0)
}
func cmdargs() {
	strout(out, config.outputCode, fmt.Sprintf(">-------Args(%v)-------<\n", len(args)-1))
	strout(out, config.outputCode, fmt.Sprintf("%q\n", args[1:]))
}
func cmdctrl() {
	b := []byte(args[1])
	x := []byte{b[0] & 0x1f}
	_, err = serialPort.Write(x)
	if err != nil {
		log.Fatal(err)
	}
	strout(out, config.outputCode, fmt.Sprintf("Ctrl+%s\n", b))
}
func cmdhex() {
	strout(out, config.outputCode, fmt.Sprintf(">-----Hex Send-----<\n"))
	strout(out, config.outputCode, fmt.Sprintf("%q\n", args[1:]))
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
	commands = append(commands, Command{name: ".ctrl", description: "发送Ctrl组合键", function: cmdctrl})
	commands = append(commands, Command{name: ".hex", description: "发送Hex", function: cmdhex})
	commands = append(commands, Command{name: ".exit", description: "退出终端", function: cmdexit})
}
