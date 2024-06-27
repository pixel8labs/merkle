package merkle

import (
	"log"

	k256 "github.com/Pixel8Labs/go-solidity-sha3"
	"github.com/ethereum/go-ethereum/common/hexutil"
	merkletree "github.com/wealdtech/go-merkletree/v2"
)

type MerkleTree struct {
	tree *merkletree.MerkleTree
}

func New(src []string) MerkleTree {
	hashType := k256.New()
	var data [][]byte
	for _, s := range src {
		data = append(data, k256.Address(s))
	}
	// Create the tree
	tree, err := merkletree.NewTree(
		merkletree.WithData(data),
		merkletree.WithHashType(hashType),
		merkletree.WithSalt(false),
		merkletree.WithSorted(true),
	)
	if err != nil {
		log.Println("failed to generate merkle tree")
	}
	// Fetch the root hash of the tree
	return MerkleTree{
		tree: tree,
	}
}

func Leaf(addr string) []byte {
	hashType := k256.New()
	return hashType.Hash(k256.Address(addr))
}

func (m MerkleTree) Root() string {
	// return crypto.Keccak256Hash(m.tree.Root()).String()
	hash := k256.New()
	return hexutil.Encode(hash.Hash(m.tree.Root()))
}

func (m MerkleTree) Proof(leaf []byte) (*merkletree.Proof, error) {
	proof, err := m.tree.GenerateProof(leaf, 0)
	if err != nil {
		return nil, err
	}

	return proof, nil
}
