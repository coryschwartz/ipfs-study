package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/subcommands"
	"github.com/multiformats/go-multihash"
)

type listCmd struct {
}

func (l *listCmd) Name() string {
	return "list"
}

func (l *listCmd) Synopsis() string {
	return "list available hash algorithms"
}

func (l *listCmd) Usage() string {
	return "list"
}

func (l *listCmd) SetFlags(f *flag.FlagSet) {
	return
}

func (c *listCmd) Execute(_ context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	for k, _ := range multihash.Names {
		fmt.Println(k)
	}
	return subcommands.ExitSuccess
}

// Perform the hash sum function on a file
type sumCmd struct {
	mhName string
}

func (c *sumCmd) Name() string {
	return "sum"
}

func (c *sumCmd) Synopsis() string {
	return "print the hash of a file"
}

func (c *sumCmd) Usage() string {
	return `sum [-h hash] <filename>:
	Print specified hashsum of file in Multihash format
	`
}

func (c *sumCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.mhName, "h", "md5", "The hash function to use")
}

func (c *sumCmd) Execute(_ context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	var filestodo []string
	if f.NArg() == 0 {
		filestodo = []string{"-"}
	} else {
		filestodo = f.Args()
	}
	for _, fn := range filestodo {
		mh, err := MultihashOfFile(fn, c.mhName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return subcommands.ExitFailure
		}
		fmt.Printf("%s  %s\n", mh.String(), fn)
	}
	return subcommands.ExitSuccess
}

func MultihashOfFile(filename string, algorithm string) (multihash.Multihash, error) {
	var buf []byte
	var err error
	if filename == "-" {
		buf, err = ioutil.ReadAll(os.Stdin)
	} else {
		buf, err = ioutil.ReadFile(filename)
	}
	if err != nil {
		return multihash.Multihash{}, err
	}
	algo, ok := multihash.Names[algorithm]
	if ok {
		return multihash.Sum(buf, algo, -1)
	} else {
		msg := fmt.Sprintf("Unsupported multihash algorithm [ %s ]", algorithm)
		return multihash.Multihash{}, errors.New(msg)
	}
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&listCmd{}, "")
	subcommands.Register(&sumCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
