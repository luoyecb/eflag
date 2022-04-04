package main

import (
	"fmt"
	"time"

	"github.com/luoyecb/eflag"
)

type CommandOptions struct {
	Name      string            `flag:"name,lycb,user name"`
	Age       int               `flag:"age,23,user age"`
	Man       bool              `flag:"man,true,user sex"`
	Salary    float64           `flag:"salary,1200.0,user salary"`
	Sleep     time.Duration     `flag:"sleep,10ms,sleep duration"`
	Addresses []string          `flag:"addr,beijing@linzhou,home address"`
	Headers   map[string]string `flag:"header,name=lisi@age=30@Content-Type=application/json,request header"`
}

func main() {
	cmdOpt := &CommandOptions{}
	eflag.Parse(cmdOpt)
	fmt.Printf("%+v\n", cmdOpt)
}
