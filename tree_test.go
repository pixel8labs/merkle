package merkle

import (
	"fmt"
	"testing"

	k256 "github.com/Pixel8Labs/go-solidity-sha3"
	"github.com/ethereum/go-ethereum/common/hexutil"
	merkletree "github.com/wealdtech/go-merkletree/v2"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		src  []string
	}{
		{"aa", []string{"0x6090A6e47849629b7245Dfa1Ca21D94cd15878Ef", "0xBE0eB53F46cd790Cd13851d5EFf43D12404d33E8"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.src)

			//Print all the nodes / leaves data
			//data := got.tree.Pollard(1)
			//for _, d := range data {
			//	fmt.Println("ERWIN DEBUG data = ", hexutil.Encode(d))
			//}

			leaf1 := k256.Address(tt.src[0])
			leaf2 := k256.Address(tt.src[1])

			// Generate a proof for address 1
			proof1, err1 := got.tree.GenerateProof(leaf1, 0)
			if err1 != nil {
				panic(err1)
			}
			// Verify the proof for leaf node
			verified1, err2 := merkletree.VerifyProofUsing(leaf1, false, proof1, [][]byte{got.tree.Root()}, k256.New())
			if err2 != nil || !verified1 {
				panic("failed to verify proof for leaf node 1")
			}
			// Generate a proof for address 2
			proof2, err3 := got.tree.GenerateProof(leaf2, 0)
			if err3 != nil {
				panic(err3)
			}
			// Verify the proof for leaf node
			verified2, err4 := merkletree.VerifyProofUsing(leaf2, false, proof2, [][]byte{got.tree.Root()}, k256.New())
			if err4 != nil || !verified2 {
				panic("failed to verify proof for leaf node 2")
			}
			// Passed, print the output
			fmt.Println("Whitelist Merkle Root = ", hexutil.Encode(got.tree.Root()))
			fmt.Println("Address = ", hexutil.Encode(leaf1), "Proof = ", hexutil.Encode(proof1.Hashes[0]))
			fmt.Println("Address = ", hexutil.Encode(leaf2), "Proof = ", hexutil.Encode(proof2.Hashes[0]))
		})
	}
}
