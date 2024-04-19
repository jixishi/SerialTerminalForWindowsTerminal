package main

type Config struct {
	portName    string
	baudRate    int
	dataBits    int
	stopBits    int
	outputCode  string
	inputCode   string
	endStr      string
	enableLog   bool
	logFilePath string
	parityBit   int
}
