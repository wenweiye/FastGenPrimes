package common

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GetRandomInt(bits int) *big.Int{
	if bits <= 0{
		panic(fmt.Errorf("GetRandomInt: bits should be positive, non-zero"))
	}
	max := new(big.Int)
	//max = max.Exp(two,big.NewInt(int64(bits)),nil).Sub(max,one)
	max = max.Exp(two,big.NewInt(int64(bits)),nil)
	n,err := IntervalRandomInt(rand.Reader, max, 1)
	if err != nil{
		panic(fmt.Errorf("rand.Int failure in GetRandomInt"))
	}
	for n.BitLen() != bits{
		n,err = IntervalRandomInt(rand.Reader, max, 1)
		if err != nil{
			panic(fmt.Errorf("rand.Int failure in GetRandomInt"))
		}
	}

	return n
}
