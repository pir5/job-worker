package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/pir5/health-worker/health_worker"
)

// @title health-worker
// @version 1.0
// @description This is PIR5 HealthCheck worker and API
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Bearer
// Commands lists the available commands and help topics.
// The order here is the order in which they are printed by 'health-worker help'.
var commands = []*health_worker.Command{
	health_worker.CmdRegister,
	health_worker.CmdWorker,
	health_worker.CmdAPI,
}

func main() {
	cmdFlags := health_worker.GlobalFlags{}
	cmdFlags.ConfPath = flag.String("config", "/etc/health-worker/worker.toml", "config file path")
	cmdFlags.PidPath = flag.String("pid", "/tmp/health-worker.pid", "pid file path")
	cmdFlags.LogPath = flag.String("logfile", "", "log file path")
	flag.Usage = usage
	flag.Parse()

	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.Flag.Usage = func() { cmd.Usage() }

			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()

			if err := cmd.Run(&cmdFlags, args); err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
			os.Exit(0)
		}
	}
	fmt.Fprintf(os.Stderr, "health-worker: unknown subcommand %q\nRun ' health-worker help' for usage.\n", args[0])
	os.Exit(2)
}

var usageTemplate = `health-worker is a tool for

Usage:

	health-worker command [arguments]

The commands are:
{{range .}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use "health-worker help [command]" for more information about a command.

`

var helpTemplate = `usage: health-worker {{.UsageLine}}

{{.Long | trim}}
`

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func printUsage(w io.Writer) {
	bw := bufio.NewWriter(w)
	tmpl(bw, usageTemplate, commands)
	bw.Flush()
}

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

// help implements the 'help' command.
func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		// not exit 2: succeeded at 'health-worker help'.
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: health-worker help command\n\nToo many arguments given.\n")
		os.Exit(2) // failed at 'health-worker help'
	}

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			// not exit 2: succeeded at 'health-worker help cmd'.
			return
		}
	}
	fmt.Fprintf(os.Stderr, "Unknown help topic %#q.  Run 'health-worker help'.\n", arg)
	os.Exit(2) // failed at 'health-worker help cmd'
}
