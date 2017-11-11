package main

import (
	"os"

	"github.com/joffotron/nepho/nepho"

	"github.com/voxelbrain/goptions"
)

func main() {
	opts := struct {
		Help goptions.Help `goptions:"-h, --help, description='Show this help'"`

		goptions.Verbs
		Create struct {
			Stack  string `goptions:"-s, --stack, description='Name of the Cloudformation stack', obligatory"`
			File   string `goptions:"-f,--file, description='Path to source single yaml file'"`
			Path   string `goptions:"-d, --dir, description='Path to source yaml files'"`
			Params string `goptions:"-p, --params, description='Parameters file'"`
		} `goptions:"create"`

		Update struct {
			Stack  string `goptions:"-s, --stack, description='Name of the Cloudformation stack', obligatory"`
			File   string `goptions:"-f,--file, description='Path to source single yaml file'"`
			Path   string `goptions:"-d, --dir, description='Path to source yaml files'"`
			Params string `goptions:"-p, --params, description='Parameters file'"`
		} `goptions:"update"`

		Diff struct {
			Stack  string `goptions:"-s, --stack, description='Name of the Cloudformation stack', obligatory"`
			File   string `goptions:"-f,--file, description='Path to source single yaml file'"`
			Path   string `goptions:"-d, --dir, description='Path to source yaml files'"`
			Params string `goptions:"-p, --params, description='Parameters file'"`
		} `goptions:"diff"`

		Delete struct {
			Stack string `goptions:"-s, --stack, description='Name of the Cloudformation stack', obligatory"`
		} `goptions:"delete"`
	}{}

	goptions.ParseAndFail(&opts)
	if len(os.Args) < 2 {
		goptions.PrintHelp()
	}

	if opts.Create.File != "" {
		if opts.Create.File != "" {
			nepho.CreateWithFile(opts.Create.Stack, opts.Create.File, opts.Create.Params)
		} else {
			nepho.CreateWithPath(opts.Create.Stack, opts.Create.Path, opts.Create.Params)
		}
	}

}
