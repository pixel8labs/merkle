package merkle

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		src  []string
		want string
	}{
		{"aa", []string{"0xD5F3603Bf3f3e673C38B5c623A4A27d20851F678", "0xa580AbDCdeAEa712Ca1E5aDa1e92A27101c64494"}, "0x328098c6dcd73990b5182359116af3bb114a90831c599394945ef5114c217478"},
		// {"aa", []string{"0xa580AbDCdeAEa712Ca1E5aDa1e92A27101c64494", "0xD5F3603Bf3f3e673C38B5c623A4A27d20851F678"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.src)

			//Print all the nodes / leaves data
			//data := got.tree.Pollard(1)
			//for _, d := range data {
			//	fmt.Println("ERWIN DEBUG data = ", hexutil.Encode(d))
			//}

			if got.Root() != tt.want {
				t.Errorf("got: %v, want: %v", got.Root(), tt.want)
				return
			}
			t.Error()
		})
	}
}
