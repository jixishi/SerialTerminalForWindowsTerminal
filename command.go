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
