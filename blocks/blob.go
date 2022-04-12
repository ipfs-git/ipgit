package blocks

import (
	"fmt"

	ipldgit "github.com/ipfs/go-ipld-git"
	"github.com/ipld/go-ipld-prime"
)

type Blob []byte

func CreateBlobNode(content Blob) (ipld.Node, error) {
	builder := ipldgit.Type.Blob.NewBuilder()

	err := builder.AssignBytes(content)
	if err != nil {
		return nil, fmt.Errorf("failed to assign content to node: %w", err)
	}

	return builder.Build(), nil
}
