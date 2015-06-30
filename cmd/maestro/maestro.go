package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hackedu/maestro/app"
	"github.com/hackedu/maestro/baton"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `maestro is maestro

Usage:

	maestro [options] command [args...]

The commands are:
	`)

		for _, c := range subcmds {
			fmt.Fprintf(os.Stderr, "    %-24s %s\n", c.name, c.description)
		}
		fmt.Fprintln(os.Stderr, `
The "maestro command -h" for more information about a command.

The options are:
`)
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
	}
	log.SetFlags(0)

	subcmd := flag.Arg(0)
	for _, c := range subcmds {
		if c.name == subcmd {
			c.run(flag.Args()[1:])
			return
		}
	}

	fmt.Fprintf(os.Stderr, "unknown subcmd %q\n", subcmd)
	fmt.Fprintln(os.Stderr, `Run "maestro -h" for usage.`)
	os.Exit(1)
}

type subcmd struct {
	name        string
	description string
	run         func(args []string)
}

var subcmds = []subcmd{
	{"serve", "start web server", serveCmd},
}

func serveCmd(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	httpAddr := fs.String("http", ":1759", "HTTP service address")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: maestro serve [options]

Starts the web server that serves the app and API.

The options are:
`)
		fs.PrintDefaults()
		os.Exit(1)
	}
	fs.Parse(args)

	if fs.NArg() != 0 {
		fs.Usage()
	}

	baton.InitModules()
	go baton.Run()
	m := http.NewServeMux()
	m.Handle("/baton/", http.StripPrefix("/baton", baton.Handler()))
	m.Handle("/", app.Handler())

	log.Print("Listening on ", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, m); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
