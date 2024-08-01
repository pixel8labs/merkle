package merkle

import (
	"crypto/sha256"
	"errors"

	"github.com/authur117/merkletree"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type MerkleTree struct {
	tree *merkletree.MerkleTree
}

func New(src []string) MerkleTree {
	var values []merkletree.Content

	for _, s := range src {
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

func (m MerkleTree) Proof(leaf string) ([][]byte, error) {
	node := merkletree.Content{
		X: leaf,
	}
	proof := m.tree.GetProof(node, 0)
	var res [][]byte
	for _, p := range proof {
		res = append(res, p.Data)
	}

	return res, nil
}

func EncodePacked(types []interface{}, values []interface{}) (string, error) {
	if len(types) != len(values) {
		return "", errors.New("params/values length mismatched")
	}

	var data []string
	for i := range types {
		_type := types[i]
		_value := values[i]
		encoded, err := Encode(_type, _value)
		if err != nil {
			return "", err
		}
		data = append(data, encoded)
	}

	return ConcatHex(data), nil
}
