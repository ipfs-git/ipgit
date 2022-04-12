package blocks

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	ipldgit "github.com/ipfs/go-ipld-git"
	"github.com/ipld/go-ipld-prime"
)

func CreateBlobNodeFromFile(path string) (ipld.Node, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %q for reading: %q", path, err)
	}
	defer file.Close()

	return CreateBlobNodeFromReader(file)
}

func CreateBlobNodeFromReader(reader io.Reader) (ipld.Node, error) {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read the content: %w", err)
	}

	builder := ipldgit.Type.Blob.NewBuilder()

	err = builder.AssignBytes(content)
	if err != nil {
		return nil, fmt.Errorf("failed to assign content to node: %w", err)
	}

	return builder.Build(), nil
}
