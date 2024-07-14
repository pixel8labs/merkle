package merkle

import (
	"crypto/sha256"
	"fmt"

	"github.com/authur117/merkletree"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type MerkleTree struct {
	tree *merkletree.MerkleTree
}

func New(src []string) MerkleTree {
	var values []merkletree.Content

	for _, s := range src {
		fmt.Println(s)
		values = append(values, merkletree.Content{
			X: s,
		})
	}

	tree := merkletree.New(merkletree.SHA3)
	tree.WithOption(merkletree.Option{
		Sort: true,
	})
	tree.InitLeaves(values)

	return MerkleTree{
		tree: tree,
	}
}

func Leaf(addr string) []byte {
	h := sha256.New()
	h.Write([]byte(addr))
	return h.Sum(nil)
}

func (m MerkleTree) Root() string {
	return hexutil.Encode(m.tree.GetRoot())
}

func (m MerkleTree) Proof(leaf merkletree.Content) ([][]byte, error) {
	proof := m.tree.GetProof(leaf, 0)
	var res [][]byte
	for _, p := range proof {
		res = append(res, p.Data)
	}

	return res, nil
}
