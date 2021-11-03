# eflag
Enhanced flag package.
Bind command line options to struct.

# example

```golang
type CommandOptions struct {
	Name      string            `eflag:"name,lycb,user name"`
	Age       int               `eflag:"age,23,user age"`
	Man       bool              `eflag:"man,true,user sex"`
	Salary    float64           `eflag:"salary,1200.0,user salary"`
	Sleep     time.Duration     `eflag:"sleep,10ms,sleep duration"`
	Addresses []string          `eflag:"addr,beijing@linzhou,home address"`
	Headers   map[string]string `eflag:"header,name=lisi@age=30,request header"`
}

func main() {
	cmdOpt := &CommandOptions{}
	eflag.Parse(cmdOpt)
	fmt.Printf("%+v\n", cmdOpt)
}

```

Test
```sh
go run eflag_demo.go -name=lisi -age=31 -man=false -salary=100 -sleep=1000000ms

# output
# &{Name:lisi Age:31 Man:false Salary:100 Sleep:16m40s}
```
