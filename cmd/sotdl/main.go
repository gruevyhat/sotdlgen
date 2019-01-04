// Package main implements a simple CLI for sotdlgen.
package main

import (
	"github.com/docopt/docopt-go"
	"github.com/gruevyhat/sotdlgen"
)

var usage = `SotDL Character Generator

Usage: sotdl [options]

Options:
  -n, --name=<str>          The character's full name; random if not specified.
  -g, --gender=<str>        The character's gender.
	-l, --level=<int>         The character's level. [default: 0]
  -A, --ancestry=<str>      The character's 0th lvl path (e.g., Human).
  -N, --novice-path=<str>   The character's 1st lvl path (e.g., Rogue). 
  -E, --expert-path=<str>   The character's 3rd lvl path (e.g., Fighter).
  -M, --master-path=<str>   The character's 7th lvl path (e.g., Myrmidon).
  -s, --seed=<hex>          Character generation signature.
  --log-level=<str>         One of {INFO, WARNING, ERROR}. [default: ERROR]
  -h --help
  --version
`

func main() {
	opts := sotdlgen.Opts{}
	optFlags, _ := docopt.ParseArgs(usage, nil, sotdlgen.VERSION)
	optFlags.Bind(&opts)
	c := sotdlgen.NewCharacter(opts)
	c.ToJSON()
}
