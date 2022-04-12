package blocks

import (
	"fmt"

	"github.com/ipfs/go-cid"
	ipldgit "github.com/ipfs/go-ipld-git"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

type TreeEntry struct {
	Mode string
	Hash datamodel.Link
}

type Tree map[string]TreeEntry

type TreeBuilder struct {
	entries map[string]TreeEntry
}

func NewTreeBuilder() *TreeBuilder {
	return &TreeBuilder{entries: map[string]TreeEntry{}}
}

func (tb *TreeBuilder) AppendLink(name string, node ipld.Node) error {

	lp := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  1,    // Usually '1'.
		Codec:    0x71, // 0x71 means "dag-cbor" -- See the multicodecs table: https://github.com/multiformats/multicodec/
		MhType:   0x13, // 0x20 means "sha2-512" -- See the multicodecs table: https://github.com/multiformats/multicodec/
		MhLength: 64,   // sha2-512 hash has a 64-byte sum.
	}}

	rawNode, err := node.AsBytes()
	if err != nil {
		return fmt.Errorf("failed to convert the node as bytes: %w", err)
	}
	lnk := lp.BuildLink(rawNode)

	tb.entries[name] = TreeEntry{
		Mode: "100644",
		Hash: lnk,
	}

	return nil
}

func (tb *TreeBuilder) Build() (ipld.Node, error) {
	builder := ipldgit.Type.Tree.NewBuilder()

	mapBuilder, err := builder.BeginMap(int64(len(tb.entries)))
	if err != nil {
		return nil, fmt.Errorf("failed to start to build the tree map: %w", err)
	}

	for key, val := range tb.entries {
		nodeBuilder, err := mapBuilder.AssembleEntry(key)
		if err != nil {
			return nil, fmt.Errorf("failed to assemble the tree entry %q: %w", key, err)
		}

		entryBuilder, err := nodeBuilder.BeginMap(2)
		if err != nil {
			return nil, fmt.Errorf("failed to start to build a TreeEntry map: %w", err)
		}

		modeBuilder, err := entryBuilder.AssembleEntry("mode")
		if err != nil {
			return nil, fmt.Errorf("failed to Assemble the entry \"mode\": %w", err)
		}
		modeBuilder.AssignString(val.Mode)

		hashBuilder, err := entryBuilder.AssembleEntry("link")
		if err != nil {
			return nil, fmt.Errorf("failed to Assemble the entry \"hash\": %w", err)
		}
		hashBuilder.AssignLink(val.Hash)
	}

	return builder.Build(), nil
}
