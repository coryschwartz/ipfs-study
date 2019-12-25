package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/subcommands"
	"github.com/multiformats/go-multihash"
)

type listCmd struct {
}

func (c *listCmd) Name() string {
	return "list"
}

func (c *listCmd) Synopsis() string {
	return "list available hash algorithms"
}

func (c *listCmd) Usage() string {
	return "list"
}

func (c *listCmd) SetFlags(f *flag.FlagSet) {
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

type checkCmd struct {
}

func (c *checkCmd) Name() string {
	return "check"
}

func (c *checkCmd) Synopsis() string {
	return "read sums from file and check them"
}

func (c *checkCmd) Usage() string {
	return `check <filename>:
	read sums from <filename> to verify files
	`
}

func (c *checkCmd) SetFlags(f *flag.FlagSet) {
	return
}

func (c *checkCmd) Execute(_ context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 1 {
		fmt.Fprintln(os.Stderr, c.Usage())
		return subcommands.ExitUsageError
	}
	sumFileName := f.Arg(0)
	sumFile, err := os.Open(sumFileName)
	if err != nil {
		panic(err)
	}
	defer sumFile.Close()
	scanner := bufio.NewScanner(sumFile)
	for scanner.Scan() {
		// similarly to an MD5SUMS file,
		// left side has the multihash string
		// right side has the filename
		fields := strings.Fields(scanner.Text())
		left, _ := multihash.FromHexString(fields[0])
		leftd, _ := multihash.Decode(left)
		right, _ := MultihashOfFile(fields[1], leftd.Name)
		if left.String() == right.String() {
			fmt.Printf("%s  PASS\n", fields[1])
		} else {
			fmt.Printf("%s  FAIL\n", fields[1])
		}
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
	subcommands.Register(&checkCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
