package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/ipfs-git/ipgit/blocks"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
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

		typ := "blob"
		if val, ok := options["type"]; ok {
			typ = val
		}

		var node ipld.Node

		switch typ {
		case "blob":
			node, err = blocks.CreateBlobNode(rawFile)
		default:
			fmt.Printf("unhandled or unimplemented type %q\n", typ)
			return 1
		}
		if err != nil {
			fmt.Printf("failed to create a blob node: %s\n", err)
			return 1
		}

		lp := cidlink.LinkPrototype{Prefix: cid.Prefix{
			Version:  1,    // Usually '1'.
			Codec:    0x71, // 0x71 means "dag-cbor" -- See the multicodecs table: https://github.com/multiformats/multicodec/
			MhType:   0x13, // 0x20 means "sha2-512" -- See the multicodecs table: https://github.com/multiformats/multicodec/
			MhLength: 64,   // sha2-512 hash has a 64-byte sum.
		}}

		rawNode, err := node.AsBytes()
		if err != nil {
			fmt.Printf("failed to convert the node as bytes: %s\n", err)
			return 1
		}
		lnk := lp.BuildLink(rawNode)

		fmt.Println(lnk)
	}

	return 0
}
