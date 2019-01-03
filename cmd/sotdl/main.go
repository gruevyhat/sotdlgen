package main

import (
	"os"

	"github.com/docopt/docopt-go"
	"github.com/gruevyhat/sotdlgen"
)

var usage = `M6IK Character Generator

Usage: m6ikgen [options]

Options:
  --name	The character's full name.
  --gender	The character's gender.
  --seed	Character generation signature.
	--log-level	One of {INFO, WARNING, ERROR}. [default: ERROR]
  -h --help
  --version
`

var Opts struct {
	Name     string `docopt:"--name"`
	Gender   string `docopt:"--gender"`
	Age      string `docopt:"--age"`
	Seed     string `docopt:"--seed"`
	LogLevel string `docopt:"--log-level"`
}

func main() {

	optFlags, _ := docopt.ParseArgs(usage, os.Args[1:], sotdlgen.VERSION)
	optFlags.Bind(&Opts)

	opts := map[string]string{
		"name":      Opts.Name,
		"gender":    Opts.Gender,
		"seed":      Opts.Seed,
		"log-level": Opts.LogLevel,
	}

	c := sotdlgen.NewCharacter(opts)
	c.Print()
}
