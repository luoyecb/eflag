package main

import (
	"fmt"
	"time"

	"github.com/luoyecb/eflag"
)

type CommandOptions struct {
	Name   string        `eflag:"name,lycb,user name"`
	Age    int           `eflag:"age,23,user age"`
	Man    bool          `eflag:"man,true,user sex"`
	Salary float64       `eflag:"salary,1200.0,user salary"`
	Sleep  time.Duration `eflag:"sleep,10ms,sleep duration"`
}

func main() {
	cmdOpt := &CommandOptions{}
	eflag.Parse(cmdOpt)
	fmt.Printf("%+v\n", cmdOpt)
}
