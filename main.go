package main

import (
	"bytes"
	"flag"
	"fmt"
	"ghubic"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"text/template"
	"unicode"
	"unicode/utf8"
)

// A Command is an implementation of a hubicli command
type Command struct {
	// Run runs the command.
	// The args are the argument after the command line.
	Run func(cmd *Command, args []string)

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'hubicli help' output.
	Short string

	// Long is the long message shown in the 'hubicli help <command>' output.
	Long string

	// Flag is the set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command will do its own
	// flag parsing.
	CustomFlags bool
}

// Name returns the command's name: the first word of the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// Usage displays the long usage of the command then exits.
func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-code.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

// This array stores all available commands
var commands = []*Command{
	cmdInfo,
    cmdAdd,
    cmdList,
}

var exitStatus = 0
var existMu sync.Mutex

func setExitStatus(n int) {
	existMu.Lock()
	if exitStatus < n {
		exitStatus = n
	}
	existMu.Unlock()
}

func printUsage(w io.Writer) {
	tmpl(w, usageTemplate, commands)
}

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

var Account *ghubic.Account

func main() {
	flag.Parse()
	flag.Usage = usage
	log.SetFlags(0)

	var err error
	Account, err = checkAuth()
	if err != nil {
		panic(err)
	}

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Runnable() {
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}
			cmd.Run(cmd, args)
			exit()
			return
		}
	}

	fmt.Fprintf(os.Stderr, "hubicli: unknown subcommand %q\nRun 'hubicli help' for usage.\n", args[0])
	setExitStatus(2)
	exit()
}

var usageTemplate = `hubicli is a command-line tool to manage your hubiC account.

Usage:

        hubicli command [arguments]

The commands are:
{{range .}}{{if .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "hubicli help [command]" for more information about a command.

Additional help topics:
{{range .}}{{if not .Runnable}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "hubicli help [topic] for more information about that topic."

`

var helpTemplate = `{{if .Runnable}}usage: hubicli {{.UsageLine}}

{{end}}{{.Long | trim}}
`

var documentationTemplate = `TODO

`

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("usage")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: hubicli help command\n\nToo many arguments given.\n")
		os.Exit(2)
	}

	arg := args[0]

	if arg == "documentation" {
		buf := new(bytes.Buffer)
		printUsage(buf)
		usage := &Command{Long: buf.String()}
		tmpl(os.Stdout, documentationTemplate, append([]*Command{usage}, commands...))
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q. Run 'hubicli help'.\n", arg)
	os.Exit(2)
}

var atexitFuncs []func()

func atexit(f func()) {
	atexitFuncs = append(atexitFuncs, f)
}

func exit() {
	for _, f := range atexitFuncs {
		f()
	}
	os.Exit(exitStatus)
}
