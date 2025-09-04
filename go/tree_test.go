package merkle

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

type leafDetail struct {
	WalletAddress string   `json:"wallet_address"`
	Cap           *big.Int `json:"cap"`
}

func TestEncodePack(t *testing.T) {
	capValue := big.NewInt(0)
	capValue.SetString("6000000000000000000", 10)

	data := []leafDetail{
		{
			WalletAddress: "0x28E5686ca4d7016281cB792CC1865D8f52444c84",
			Cap:           capValue,
		},
	}

	t.Run("EncodePacked test", func(t *testing.T) {
		var encodePackedResult []string
		for _, d := range data {
			res, err := EncodePacked(
				[]interface{}{
					"address", "string", "uint256",
				},
				[]interface{}{
					d.WalletAddress, ",", d.Cap,
				})
			if err != nil {
				t.Fatalf("EncodePacked failed: %v", err)
			}
			encodePackedResult = append(encodePackedResult, res)
		}

		// Check if the results are as expected
		expectedResults := []string{
			// Update these expected results with actual expected encoded values
			"0x28e5686ca4d7016281cb792cc1865d8f52444c842c00000000000000000000000000000000000000000000000053444835ec580000",
		}

		if encodePackedResult[0] != expectedResults[0] {
			t.Errorf("got: %v, want: %v", encodePackedResult[0], expectedResults[0])
			return
		}
	})
}

func TestNew(t *testing.T) {
	capValue := big.NewInt(0)
	capValue.SetString("20000000000000000000", 10)

	capValue2 := big.NewInt(0)
	capValue2.SetString("5000000000000000000", 10)

	data := []leafDetail{
		{
			WalletAddress: "0xd5f3603bf3f3e673c38b5c623a4a27d20851f678",
			Cap:           capValue,
		},
		{
			WalletAddress: "0xB6ecE3bDC66810f3CbC4697A8baC9A37E300ef29",
			Cap:           capValue2,
		},
	}

	var encodePackedResult []string
	for _, d := range data {
		res, err := EncodePacked(
			[]interface{}{
				"address", "string", "uint256",
			},
			[]interface{}{
				d.WalletAddress, ",", d.Cap,
			})
		if err != nil {
			fmt.Println(err)
			return
		}
		encodePackedResult = append(encodePackedResult, res)
	}

	tests := []struct {
		name string
		src  []string
		want string
	}{
		{"aa", encodePackedResult, "0x2e23f305b02010697b5769b3e5d1895ee2019880c43072b4aa12cbc6597ff7e4"},
		// 0x160fa9c92d0a4fd6cee96e18b7eb6e041e2b2c1657d21f6913f72c8265eaddb7
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.src)

			if got.Root() != tt.want {
				t.Errorf("got: %v, want: %v", got.Root(), tt.want)
				return
			}
			// t.Error()
		})
	}
}

func TestProof_ExistingLeaf(t *testing.T) {
	// Arrange
	leaves := []string{"0x1b094c29db41629a080553194abb192f6213c0162c0000000000000000000000000000000000000000000000008ac7230489e80000", "0xada8bf34b136194fc9a95d089d8eedf28b1b19562c0000000000000000000000000000000000000000000000008ac7230489e80000"}
	mt := New(leaves)

	proof, err := mt.Proof("0x1b094c29db41629a080553194abb192f6213c0162c0000000000000000000000000000000000000000000000008ac7230489e80000")
	if err != nil {
		t.Fatalf("Proof returned error: %v", err)
	}

	// Assert: proof should not be empty for an existing leaf
	if len(proof) == 0 {
		t.Fatalf("expected non-empty proof for existing leaf, got len=0")
	}

	for i, sib := range proof {
		t.Logf("proof[%d]: 0x%s", i, hex.EncodeToString(sib))
		fmt.Printf("proof[%d]: 0x%s", i, hex.EncodeToString(sib))
	}

	// Expected Proof: 0x2f79c9672e91eb7527c1d564f38b129e14345f75f564981b61ba40a3203b150d

	// Each proof node should be a 32-byte hash (most merkle libs use 32 bytes)
	for i, sib := range proof {
		if len(sib) == 0 {
			t.Fatalf("proof[%d] was empty", i)
		}
		// Typically 32 bytes; if your library uses a different digest length, adjust here.
		if len(sib) != 32 {
			t.Fatalf("proof[%d] unexpected length: got %d, want 32", i, len(sib))
		}
	}
}
