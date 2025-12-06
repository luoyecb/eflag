# eflag
Enhanced std flag package.
Bind command-line options to struct.

# example

```golang
type CommandOptions struct {
	Name      string            `flag:"name" flag_short:"n" default:"lycb" usage:"user name"`
	Age       int               `flag:"age" default:"23" usage:"user age"`
	Man       bool              `flag:"man" default:"false" usage:"user sex"`
	Salary    float64           `flag:"salary" default:"1200.0" usage:"user salary"`
	Sleep     time.Duration     `flag:"sleep" default:"10ms" usage:"sleep duration"`
	Addresses []string          `flag:"addr" default:"beijing@linzhou" usage:"home address"`
	Headers   map[string]string `flag:"header" default:"name=lisi@age=30@Content-Type=application/json" usage:"request header"`

	ShowList   bool   `flag:"show_list" default:"false" usage:"show list" command:""`
	ShowDetail bool   `flag:"show_detail" default:"true" usage:"show detail" command:",false"`
	Cover      string `flag:"cover" default:"" usage:"add cover" command:""`

	Show   string `sub_command:"show" usage:"show action"`
	Detail bool   `sub_command:"detail" usage:"show detail action"`

	Args []string
}

func (opt *CommandOptions) HeadersDefault() map[string]string {
	return map[string]string{
		"lang": "golang",
	}
}

func (opt *CommandOptions) ManDefault() bool {
	return true
}

func (opt *CommandOptions) ShowCommand() {
	fmt.Println("show sub_command")
}

func (opt *CommandOptions) ShowListCommand() {
	fmt.Println("show list")
}

func (opt *CommandOptions) ShowDetailCommand() {
	fmt.Println("show detail")
}

func (opt *CommandOptions) CoverCommand() {
	fmt.Println("cover")
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
