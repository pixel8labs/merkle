package merkle

import (
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
	capValue.SetString("6000000000000000000", 10)

	capValue2 := big.NewInt(0)
	capValue2.SetString("9000000000000000000", 10)

	data := []leafDetail{
		{
			WalletAddress: "0x28E5686ca4d7016281cB792CC1865D8f52444c84",
			Cap:           capValue,
		},
		{
			WalletAddress: "0x1548AaE0d0eC35d31a199FBA1A4F56282DE73A15",
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
