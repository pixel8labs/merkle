package merkle

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type PadOptions struct {
	Dir  string // "right" or "left"
	Size int
}

type SizeExceedsPaddingSizeError struct {
	Size       int
	TargetSize int
	Type       string
}

type NumberToHexOpts struct {
	Signed bool
	Size   int
}

func ConcatHex(data []string) string {
	if len(data) == 0 {
		return "0x"
	}

	var sb strings.Builder
	sb.WriteString("0x")
	for _, h := range data {
		sb.WriteString(h[2:]) // Remove "0x" prefix and concatenate
	}
	return sb.String()
}

func Encode[T any](valueType T, value T) (string, error) {
	if any(valueType) == Address {
		// It's more like a validation because we should accept 0x{string} address
		address, ok := any(value).(string)
		if !ok {
			return "", errors.New("value is not a string")
		}
		if !common.IsHexAddress(address) {
			return "", errors.New("address is not valid")
		}
		padHex, err := PadHex(strings.ToLower(address), PadOptions{
			Size: 0,
			Dir:  "left",
		})
		if err != nil {
			return "", err
		}
		if !isValidHex(padHex) {
			return "", errors.New("not a valid hex string")
		}
		return padHex, nil
	}

	if any(valueType) == String {
		str, ok := any(value).(string)
		if !ok {
			return "", errors.New("value is not a string")
		}
		return hexutil.Encode([]byte(str)), nil
	}

	if any(valueType) == Uint256 {
		size := 32
		hexResult, err := NumberToHex(value, NumberToHexOpts{
			Signed: false,
			Size:   size,
		})
		if err != nil {
			return "", err
		}
		return hexResult, nil
	}

	return "", nil
}

func NumberToHex[T any](value T, opts NumberToHexOpts) (string, error) {
	var bigValue *big.Int

	switch v := any(value).(type) {
	case int:
		bigValue = big.NewInt(int64(v))
	case int64:
		bigValue = big.NewInt(v)
	case uint64:
		bigValue = new(big.Int).SetUint64(v)
	case *big.Int:
		bigValue = new(big.Int).Set(v)
	default:
		return "", errors.New("unsupported value type")
	}

	var maxValue, minValue *big.Int
	maxValue = new(big.Int).SetUint64(uint64(math.MaxInt64)) // Maximum value for int
	minValue = big.NewInt(0)

	if opts.Size > 0 {
		bits := uint(opts.Size * 8)
		maxValue = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), bits), big.NewInt(1))
		if opts.Signed {
			maxValue = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), bits-1), big.NewInt(1))
			minValue = new(big.Int).Neg(new(big.Int).Add(maxValue, big.NewInt(1)))
		}
	}

	if (maxValue != nil && bigValue.Cmp(maxValue) > 0) || bigValue.Cmp(minValue) < 0 {
		return "", fmt.Errorf("integer %s out of range (min: %s, max: %s, signed: %v, size: %d)",
			bigValue.String(),
			minValue.String(),
			maxValue.String(),
			opts.Signed,
			opts.Size)
	}

	hexStr := bigValue.Text(16)
	if opts.Size > 0 {
		return PadHex(hexStr, PadOptions{
			Size: opts.Size,
			Dir:  "left",
		})
	}

	return "0x" + hexStr, nil
}

func PadHex(hex string, options PadOptions) (string, error) {
	if options.Size == 0 {
		return hex, nil
	}

	hexStr := strings.TrimPrefix(string(hex), "0x")
	if len(hexStr) > options.Size*2 {
		return "", fmt.Errorf("Size (%d) exceeds padding size (%d)", len(hexStr)/2, options.Size)
	}

	padSize := options.Size * 2
	if options.Dir == "right" {
		hexStr = hexStr + strings.Repeat("0", padSize-len(hexStr))
	} else {
		hexStr = strings.Repeat("0", padSize-len(hexStr)) + hexStr
	}

	return "0x" + hexStr, nil
}

func isValidHex(address string) bool {
	re := regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	return re.MatchString(address)
}
