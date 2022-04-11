package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/teris-io/cli"
)

func HashObjectCmd() cli.Command {
	return cli.NewCommand("hash-object",
		"Computes the object ID value for an object with specified type with the contents of the named file (which can be outside of the work tree), and optionally writes the resulting object into the object database. Reports its object ID to its standard output. When <type> is not specified, it defaults to \"blob\".").
		WithArg(cli.NewArg("file...", "List of file to hash")).
		WithOption(cli.NewOption("type", "Specify the type (default: \"blob\").").WithChar('t').WithType(cli.TypeString)).
		WithAction(hashObjectAction)
}

func hashObjectAction(args []string, options map[string]string) int {
	for _, arg := range args {
		rawFile, err := ioutil.ReadFile(arg)
		if err != nil {
			fmt.Printf("fatal: could not open %q for reading: %s\n", arg, err)
			return 1
		}

		blockType := "blob"
		if val, ok := options["type"]; ok {
			blockType = val
		}

		fmt.Printf("%s => %s", blockType, rawFile)
	}

	return 0
}
