package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)
const VERSION string = "v0.11"
const DATESTRING string = "june 2025"
const HELP_MSG string = 
`Usage:
  txtr <input> <output> [options]

Options:
  --kvinput      Input for a keyvalue file

Examples:
  txtr docker-compose.example.yml docker-compose.yml --kvinput .env`

func Help() {
	Version()
	fmt.Println(HELP_MSG)
}
func Version() {
	fmt.Printf("txtr version %v by jack anderson, %v\n", VERSION, DATESTRING)
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
	aliases []string
	handle  func(Option, *Command, []byte) ([]byte, error)
}

// Dirty copy of the fmt.printf function, wrapped to only print on verbose
func (c Command) Vlog(format string, a ...any) (n int, err error) {
	if !c.Flag("v") {
		return 0, nil
	}
	return fmt.Fprintf(os.Stdout, format, a...)
}

type Command struct {
	args           Args
	bytes          []byte
	flags          []string
	options        []Option
	optionHandlers []OptHandler
}

func (c Command) Flag(key string) bool {
	return slices.Contains(c.flags, key)
}
func (c Command) GetOption(alias string) []Option {
	var opts []Option
	for i := range c.options {
		opt := c.options[i]
		if opt.name == alias {
			return append(opts, opt)
		}
	}
	return nil
}
func (c Command) RunOpts() error {
	for i := range c.options {
		opt := c.options[i]

		for hi := range c.optionHandlers {
			h := c.optionHandlers[hi]
			if slices.Contains(h.aliases, opt.name) {
				b, err := h.handle(opt, &c, c.bytes)
				if err != nil {
					return err
				}
				c.bytes = b
			}
		}
	}
	return nil
}
func (c Command) ScanFlags(args []string) []string {
	var flags []string
	for i := range args {
		arg := args[i]
		chars := []rune(arg)
		if string(chars[0]) == "-" && string(chars[1]) != "-" {
			chars := string(chars[1:])
			c.Vlog("Found flags `%v`\n", chars)
			flags = append(flags, strings.Split(chars, "")...)
		}
	}
	return flags
}
func (c Command) ScanOpts(args []string, starti int, options []Option) []Option {
	var optname []string
	var optargs []string
	for i := starti; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") && len(optname) == 0 {
			optname = append(optname, strings.Replace(arg, "--", "", 1))
			c.Vlog("Found option `%v`\n", arg)
			continue
		} else if strings.HasPrefix(arg, "--") && len(optname) > 0 {
			c.Vlog("Found option `%v`\n", arg)
			c.Vlog("Committing option `%v`\n", optname[0])
			// Commit the current option, then move on
			options = append(options, Option{
				name: optname[0],
				args: optargs,
			})
			return c.ScanOpts(args, i, options)
		} else if !strings.HasPrefix(arg, "-") {
			c.Vlog("Found argument `%v` belonging to option `%v`\n", arg, optname[0])
			optargs = append(optargs, arg)
		}
	}
	// Commit the current option if there's not one after it, 
	// and if there is an option to begin with
	if len(optname) > 0 {
		c.Vlog("Committing option `%v`\n", optname[0])
		options = append(options, Option{
			name: optname[0],
			args: optargs,
		})
	}

	return options
}
func ParseCommand(args []string, optHandlers []OptHandler) (*Command, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments")
	}
	c := Command{optionHandlers: optHandlers}
	c.args.input = args[0]
	c.args.output = args[1]

	ib, err := os.ReadFile(c.args.input)
	if err != nil {
		return nil, err
	}
	c.bytes = ib

	c.flags = c.ScanFlags(args[2:])
	c.options = c.ScanOpts(args[2:], 0, make([]Option, 0))
	return &c, nil
}

func main() {
	var optHandlers []OptHandler
	optHandlers = append(optHandlers,
		OptHandler{aliases: append(make([]string, 0), "kvinput"), handle: Kv_Run},
	)

	// Provide only useful args
	cmd, err := ParseCommand(os.Args[1:], optHandlers)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		Help()
		return
	}
	if err := cmd.RunOpts(); err != nil {
		fmt.Printf("Error: %v\n", err)
		Help()
		return
	}
	os.WriteFile(cmd.args.output, cmd.bytes, os.ModePerm)
}
