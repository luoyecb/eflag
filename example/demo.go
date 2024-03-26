package main

import (
	"fmt"
	"time"

	"github.com/luoyecb/eflag"
)

type CommandOptions struct {
	Name      string            `flag:"name" default:"lycb" usage:"user name"`
	Age       int               `flag:"age" default:"23" usage:"user age"`
	Man       bool              `flag:"man" default:"false" usage:"user sex"`
	Salary    float64           `flag:"salary" default:"1200.0" usage:"user salary"`
	Sleep     time.Duration     `flag:"sleep" default:"10ms" usage:"sleep duration"`
	Addresses []string          `flag:"addr" default:"beijing@linzhou" usage:"home address"`
	Headers   map[string]string `flag:"header" default:"name=lisi@age=30@Content-Type=application/json" usage:"request header"`

	ShowList   bool `flag:"show_list" default:"false" usage:"show list" command:""`
	ShowDetail bool `flag:"show_detail" default:"false" usage:"show detail" command:""`
}

func (opt *CommandOptions) HeadersDefault() map[string]string {
	return map[string]string{
		"lang": "golang",
	}
}

func (opt *CommandOptions) ShowListCommand() {
	fmt.Println("show list")
}

func (opt *CommandOptions) ShowDetailCommand() {
	fmt.Println("show detail")
}

func main() {
	cmdOpt := &CommandOptions{}

	// eflag.Parse(cmdOpt)
	// fmt.Printf("%+v\n", cmdOpt)
	// eflag.RunCommand(cmdOpt)

	eflag.ParseAndRunCommand(cmdOpt)
	fmt.Printf("%+v\n", cmdOpt)
}
