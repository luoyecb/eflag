# eflag
Enhanced std flag package.
Bind command-line options to struct.

# example

```golang
type CommandOptions struct {
	Name      string            `flag:"name" default:"lycb" usage:"user name"`
	Age       int               `flag:"age" default:"23" usage:"user age"`
	Man       bool              `flag:"man" default:"true" usage:"user sex"`
	Salary    float64           `flag:"salary" default:"1200.0" usage:"user salary"`
	Sleep     time.Duration     `flag:"sleep" default:"10ms" usage:"sleep duration"`
	Addresses []string          `flag:"addr" default:"beijing@linzhou" usage:"home address"`
	Headers   map[string]string `flag:"header" default:"name=lisi@age=30@Content-Type=application/json" usage:"request header"`
}

func (opt *CommandOptions) HeadersDefault() map[string]string {
	return map[string]string{
		"lang": "golang",
	}
}

func main() {
	cmdOpt := &CommandOptions{}
	eflag.Parse(cmdOpt)
	fmt.Printf("%+v\n", cmdOpt)
}

```

Test
```sh
go run example/demo.go -name=lisi -age=31 -man=false -salary=100 -sleep=1000000ms

# output
# &{Name:lisi Age:31 Man:false Salary:100 Sleep:16m40s Addresses:[beijing linzhou] Headers:map[lang:golang]}
```
