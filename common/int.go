package common

import (
	"io"
	"math/big"
)

type modInt big.Int

var (
	zero = big.NewInt(0)
	one = big.NewInt(1)
	two = big.NewInt(2)
	five = big.NewInt(5)
)


func IntervalRandomInt(rand io.Reader, max *big.Int, tag int) (n *big.Int, err error) {
	if max.Sign() <= 0 {
		panic("crypto/rand: argument to Int is <= 0")
	}
	n = new(big.Int)
	n.Sub(max, n.SetUint64(1))
	// bitLen is the maximum bit length needed to encode a value < max.
	bitLen := n.BitLen()
	if bitLen == 0 {
		// the only valid result is 0
		return
	}
	// k is the maximum byte length needed to encode a value < max.
	k := (bitLen + 7) / 8
	// b is the number of bits in the most significant byte of max-1.
	b := uint(bitLen % 8)
	if b == 0 {
		b = 8
	}

	bytes := make([]byte, k)

	for {
		_, err = io.ReadFull(rand, bytes)
		if err != nil {
			return nil, err
		}

		if tag != -1{
			bytes[k-1] |= byte(tag)
		}
		// Clear bits in the first byte to increase the probability
		// that the candidate is < max.
		bytes[0] &= uint8(int(1<<b) - 1)
		bytes[0] |= 128
		n.SetBytes(bytes)
		if n.Cmp(max) < 0 {
			return
		}
	}
}