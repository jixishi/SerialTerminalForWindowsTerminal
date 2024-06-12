package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	inf "github.com/fzdwx/infinite"
	"github.com/fzdwx/infinite/color"
	"github.com/fzdwx/infinite/components"
	"github.com/fzdwx/infinite/components/input/text"
	"github.com/fzdwx/infinite/components/selection/confirm"
	"github.com/fzdwx/infinite/components/selection/singleselect"
	"github.com/fzdwx/infinite/style"
	"github.com/fzdwx/infinite/theme"
	"github.com/spf13/pflag"
	"go.bug.st/serial"
	"log"
	"strconv"
	"strings"
)

type ptrVal struct {
	*string
	sl *[]string
	*int
	il *[]int
	*bool
	*float64
	*float32
	ext *string
}
type Val struct {
	string
	int
	bool
	float64
	float32
	extdef string
}
type Flag struct {
	v    ptrVal
	sStr string
	lStr string
	dv   Val
	help string
}

var (
	portName   = Flag{ptrVal{string: &config.portName}, "p", "port", Val{string: ""}, "要连接的串口\t(/dev/ttyUSB0、COMx)"}
	baudRate   = Flag{ptrVal{int: &config.baudRate}, "b", "baud", Val{int: 115200}, "波特率"}
	dataBits   = Flag{ptrVal{int: &config.dataBits}, "d", "data", Val{int: 8}, "数据位"}
	stopBits   = Flag{ptrVal{int: &config.stopBits}, "s", "stop", Val{int: 0}, "停止位停止位(0: 1停止 1:1.5停止 2:2停止)"}
	outputCode = Flag{ptrVal{string: &config.outputCode}, "o", "out", Val{string: "UTF-8"}, "输出编码"}
	inputCode  = Flag{ptrVal{string: &config.inputCode}, "i", "in", Val{string: "UTF-8"}, "输入编码"}
	endStr     = Flag{ptrVal{string: &config.endStr}, "e", "end", Val{string: "\n"}, "终端换行符"}
	logExt     = Flag{v: ptrVal{ext: &config.logFilePath}, sStr: "l", lStr: "log", dv: Val{extdef: "./%s-$s.txt", string: ""}, help: "日志保存路径"}
	timeExt    = Flag{v: ptrVal{ext: &config.timesFmt}, sStr: "t", lStr: "time", dv: Val{extdef: "[06-01-02 15:04:05.000]", string: ""}, help: "时间戳格式化字段"}
	forWard    = Flag{ptrVal{il: &config.forWard}, "f", "forward", Val{int: 0}, "转发模式(0: 无 1:TCP-C 2:UDP-C 支持多次传入)"}
	address    = Flag{ptrVal{sl: &config.address}, "a", "address", Val{string: "127.0.0.1:12345"}, "转发服务地址(支持多次传入)"}
	frameSize  = Flag{ptrVal{int: &config.frameSize}, "F", "Frame", Val{int: 16}, "帧大小"}
	parityBit  = Flag{ptrVal{int: &config.parityBit}, "v", "verify", Val{int: 0}, "奇偶校验(0:无校验、1:奇校验、2:偶校验、3:1校验、4:0校验)"}
	flags      = []Flag{portName, baudRate, dataBits, stopBits, outputCode, inputCode, endStr, forWard, address, frameSize, parityBit, logExt, timeExt}
)

var (
	bauds = []string{"自定义", "300", "600", "1200", "2400", "4800", "9600",
		"14400", "19200", "38400", "56000", "57600", "115200", "128000",
		"256000", "460800", "512000", "750000", "921600", "1500000"}
	datas    = []string{"5", "6", "7", "8"}
	stops    = []string{"1", "1.5", "2"}
	paritys  = []string{"无校验", "奇校验", "偶校验", "1校验", "0校验"}
	forwards = []string{"No", "TCP-C", "UDP-C"}
)

type ValType int

const (
	notVal ValType = iota
	stringVal
	intVal
	boolVal
	extVal
)

func printUsage(ports []string) {
	fmt.Printf("\n参数帮助:\n")
	for _, f := range flags {
		flagprint(f)
	}
	fmt.Printf("\n在线串口: %v\n", strings.Join(ports, ","))
}
func flagFindValue(v ptrVal) ValType {
	if v.string != nil {
		return stringVal
	}
	if v.bool != nil {
		return boolVal
	}
	if v.int != nil {
		return intVal
	}
	if v.ext != nil {
		return extVal
	}
	return notVal
}
func flagprint(f Flag) {
	switch flagFindValue(f.v) {
	case stringVal:
		fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%q\n", f.sStr, f.lStr, f.dv.string, f.help, f.dv.string)
	case intVal:
		fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", f.sStr, f.lStr, f.dv.int, f.help, f.dv.int)
	case boolVal:
		fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", f.sStr, f.lStr, f.dv.bool, f.help, f.dv.bool)
	case extVal:
		fmt.Printf("\t-%v -%v %T \n\t  %v\t默认值:%v\n", f.sStr, f.lStr, f.dv.extdef, f.help, f.dv.extdef)
	default:
		panic("unhandled default case")
	}
}
func flagInit(f *Flag) {
	if f.v.string != nil {
		pflag.StringVarP(f.v.string, f.lStr, f.sStr, f.dv.string, f.help)
	}
	if f.v.bool != nil {
		pflag.BoolVarP(f.v.bool, f.lStr, f.sStr, f.dv.bool, f.help)
	}
	if f.v.int != nil {
		pflag.IntVarP(f.v.int, f.lStr, f.sStr, f.dv.int, f.help)
	}
	if f.v.ext != nil {
		pflag.StringVarP(f.v.ext, f.lStr, f.sStr, f.dv.string, f.help)
		pflag.Lookup(f.lStr).NoOptDefVal = f.dv.extdef
	}
	if f.v.sl != nil {
		pflag.StringArrayVarP(f.v.sl, f.lStr, f.sStr, []string{f.dv.string}, f.help)
	}
	if f.v.il != nil {
		pflag.IntSliceVarP(f.v.il, f.lStr, f.sStr, []int{f.dv.int}, f.help)
	}
}
func flagExt() {
	if config.logFilePath != "" {
		config.enableLog = true
	}
	if config.timesFmt != "" {
		config.timesTamp = true
	}
}
func getCliFlag() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}

	inputs := components.NewInput()
	inputs.Prompt = "Filtering: "
	inputs.PromptStyle = style.New().Bold().Italic().Fg(color.LightBlue)

	selectKeymap := singleselect.DefaultSingleKeyMap()
	selectKeymap.Confirm = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "finish select"),
	)
	selectKeymap.Choice = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "finish select"),
	)
	selectKeymap.NextPage = key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("->", "next page"),
	)
	selectKeymap.PrevPage = key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("<-", "prev page"),
	)

	s, _ := inf.NewSingleSelect(
		ports,
		singleselect.WithKeyBinding(selectKeymap),
		singleselect.WithPageSize(4),
		singleselect.WithFilterInput(inputs),
	).Display("选择串口")
	config.portName = ports[s]

	s, _ = inf.NewSingleSelect(
		bauds,
		singleselect.WithKeyBinding(selectKeymap),
		singleselect.WithPageSize(4),
	).Display("选择波特率")
	if s != 0 {
		config.baudRate, _ = strconv.Atoi(bauds[s])
	} else {
		b, _ := inf.NewText(
			text.WithPrompt("BaudRate:"),
			text.WithPromptStyle(theme.DefaultTheme.PromptStyle),
			text.WithDefaultValue("115200"),
		).Display()
		config.baudRate, _ = strconv.Atoi(b)
	}
	v, _ := inf.NewConfirmWithSelection(
		confirm.WithPrompt("启用Hex"),
	).Display()
	if v {
		config.inputCode = "hex"
		b, _ := inf.NewText(
			text.WithPrompt("Frames:"),
			text.WithPromptStyle(theme.DefaultTheme.PromptStyle),
			text.WithDefaultValue("16"),
		).Display()
		config.frameSize, _ = strconv.Atoi(b)
	}
	v, _ = inf.NewConfirmWithSelection(
		confirm.WithPrompt("启用时间戳"),
	).Display()
	config.timesTamp = v
	if v {
		b, _ := inf.NewText(
			text.WithPrompt("格式化字段:"),
			text.WithPromptStyle(theme.DefaultTheme.PromptStyle),
			text.WithDefaultValue(timeExt.dv.extdef),
		).Display()
		config.timesFmt = b
	}
	v, _ = inf.NewConfirmWithSelection(
		confirm.WithPrompt("启用高级配置"),
	).Display()
	if v {
		s, _ = inf.NewSingleSelect(
			datas,
			singleselect.WithKeyBinding(selectKeymap),
			singleselect.WithPageSize(4),
			singleselect.WithFilterInput(inputs),
		).Display("选择数据位")
		config.dataBits, _ = strconv.Atoi(datas[s])

		s, _ = inf.NewSingleSelect(
			stops,
			singleselect.WithKeyBinding(selectKeymap),
			singleselect.WithPageSize(4),
			singleselect.WithFilterInput(inputs),
		).Display("选择停止位")
		config.stopBits = s

		s, _ = inf.NewSingleSelect(
			paritys,
			singleselect.WithKeyBinding(selectKeymap),
			singleselect.WithPageSize(4),
			singleselect.WithFilterInput(inputs),
		).Display("选择校验位")
		config.parityBit = s

		t, _ := inf.NewText(
			text.WithPrompt("换行符:"),
			text.WithPromptStyle(theme.DefaultTheme.PromptStyle),
			text.WithDefaultValue(endStr.dv.string),
		).Display()
		config.endStr = t

		v, _ = inf.NewConfirmWithSelection(
			confirm.WithDefaultYes(),
			confirm.WithPrompt("启用编码转换"),
		).Display()

		if v {
			t, _ = inf.NewText(
				text.WithPrompt("输入编码:"),
				text.WithPromptStyle(theme.DefaultTheme.PromptStyle),
				text.WithDefaultValue(inputCode.dv.string),
			).Display()
			config.inputCode = t

			t, _ = inf.NewText(
				text.WithPrompt("输出编码:"),
				text.WithPromptStyle(theme.DefaultTheme.PromptStyle),
				text.WithDefaultValue(outputCode.dv.string),
			).Display()
			config.outputCode = t
		}
	G_F_mode:
		s, _ = inf.NewSingleSelect(
			forwards,
			singleselect.WithKeyBinding(selectKeymap),
			singleselect.WithPageSize(3),
			singleselect.WithFilterInput(inputs),
		).Display("选择转发模式")
		if s != 0 {
			config.forWard = append(config.forWard, s)
			t, _ = inf.NewText(
				text.WithPrompt("地址:"),
				text.WithPromptStyle(theme.DefaultTheme.PromptStyle),
				text.WithDefaultValue(address.dv.string),
			).Display()
			config.address = append(config.address, t)
			goto G_F_mode
		}

		e, _ := inf.NewConfirmWithSelection(
			confirm.WithDefaultYes(),
			confirm.WithPrompt("启用日志"),
		).Display()
		config.enableLog = e
		if e {
			t, _ = inf.NewText(
				text.WithPrompt("Path:"),
				text.WithPromptStyle(theme.DefaultTheme.PromptStyle),
				text.WithDefaultValue("./%s-$s.txt"),
			).Display()
			config.logFilePath = t
		}
	}

}
