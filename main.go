package main

import (
	"fmt"
	"os"
	"strings"
)

const HELP_MSG string = `Usage:
  txtr <input> <output> [options]

Options:
  --kvinput      Input for a keyvalue file

Examples:
  txtr docker-compose.example.yml docker-compose.yml --kvinput .env`

func Help() {
	fmt.Println(HELP_MSG)
}

type Args struct {
	input  string
	output string
}
type Option struct {
	name string
	args []string
}
type OptHandler struct {
	name string
	handle func(Option, []byte) ([]byte, error)
}

type Command struct {
	args    Args
	options []Option
	optionHandlers []OptHandler
}

func (c Command) RunOpts() ([]byte, error) {
	proc_bytes, err := os.ReadFile(c.args.input)
	if err != nil {
		return nil, err
	}
	for i := range c.options {
		opt := c.options[i]
		
		for hi := range c.optionHandlers {
			h := c.optionHandlers[hi]
			if h.name == opt.name {
				b, err := h.handle(opt, proc_bytes)
				if err != nil {
					return nil, err
				}
				proc_bytes = b
			}
		}
	}
	return proc_bytes, nil
}
func ScanOpts(args []string, starti int, options []Option) []Option {
	var optname []string
	var optargs []string
	for i := starti; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") && len(optname) == 0 {
			optname = append(optname, strings.Replace(arg, "--", "", 1))
			fmt.Printf("Found option `%v`\n", arg)
			continue
		} else if strings.HasPrefix(arg, "--") && len(optname) > 0 {
			fmt.Printf("Found option `%v`\n", arg)
			fmt.Printf("Committing option `%v`\n", optname[0])
			// Commit the current option, then move on
			options = append(options, Option{
				name: optname[0],
				args: optargs,
			})
			return ScanOpts(args, i, options)
		} else {
			fmt.Printf("Found argument `%v` belonging to option `%v`\n", arg, optname[0])
			optargs = append(optargs, arg)
		}
	}
	// Commit the current option if there's not one after it
	fmt.Printf("Committing option `%v`\n", optname[0])
	options = append(options, Option{
		name: optname[0],
		args: optargs,
	})
	return options
}
func NewCommand(args []string, optHandlers []OptHandler) *Command {
	if len(args) < 2 {
		return nil
	}
	c := Command{optionHandlers: optHandlers}
	c.args.input = args[0]
	c.args.output = args[1]

	opts := ScanOpts(args[2:], 0, make([]Option, 0))
	c.options = opts
	return &c
}

func main() {
	var optHandlers []OptHandler;
	optHandlers = append(optHandlers, 
		OptHandler{name: "kvinput", handle: Kv_Run},
	)

	// Provide only useful args
	cmd := NewCommand(os.Args[1:], optHandlers)
	if cmd == nil {
		Help()
		return
	}
	output, err := cmd.RunOpts()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		Help()
		return
	}
	os.WriteFile(cmd.args.output, output, os.ModePerm)
}


