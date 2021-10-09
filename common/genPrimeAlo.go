package common

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

func GenCrtM(len int) *big.Int {
	n := big.NewInt(1)
	for i := 0; i < len; i++ {
		n = n.Mul(n, big.NewInt(int64(primes[i])))
	}
	return n
}

func TraPrimeGen(times, bits int) {
	start := time.Now()
	count := 0
	primesCount := 1
	p := new(big.Int)
	p = GetRandomInt(bits)
	for count < times {
		for !p.ProbablyPrime(30) {
			primesCount++
			p = GetRandomInt(bits)
		}
		p = GetRandomInt(bits)
		count++
	}
	end := time.Since(start)
	fmt.Printf("传统素数平均生成时间%fms\n", float64(end.Milliseconds())/float64(count))
	fmt.Printf("传统平均素性测试次数%d\n", primesCount/count)
}

func TraPrimeGenWithIncreasing(times, bits int) {
	start := time.Now()
	count := 0
	primesCount := 1
	p := new(big.Int)
	p = GetRandomInt(bits)
	for count < times {
		for !p.ProbablyPrime(30) {
			primesCount++
			p.Add(p, two)
			if p.BitLen() > bits {
				p = GetRandomInt(bits)
			}
		}
		p = GetRandomInt(bits)
		count++
	}
	end := time.Since(start)
	fmt.Printf("增量素数平均生成时间%fms\n", float64(end.Milliseconds())/float64(count))
	fmt.Printf("增量平均素性测试次数%d\n", primesCount/count)
}

func TraPrimeGenWithImProIncreasing(times, bits, len int) {
	start := time.Now()
	count := 0
	primesCount := 1
	pi := GenCrtM(len)

	p, q := new(big.Int), new(big.Int)

	for count < times {
	refresh:
		p = GetRandomInt(bits)
		q.GCD(big.NewInt(0), big.NewInt(0), p, pi)
		for !bool(q.Cmp(one) == 0) {
			p = GetRandomInt(bits)
			q.GCD(big.NewInt(0), big.NewInt(0), p, pi)
		}

		for !p.ProbablyPrime(30) {
			primesCount++
			p = p.Add(p, pi)

			if p.BitLen() > bits {
				goto refresh
			}
		}
		count++
	}
	end := time.Since(start)
	fmt.Printf("改进的增量素数平均生成时间%fms\n", float64(end.Milliseconds())/float64(count))
	fmt.Printf("改进的增量素性测试次数%d\n", primesCount/count)
}

func MJSpecialGenPrimes(times, bits, len int) {
	omegaMin, omegaMax := new(big.Int), new(big.Int)
	omegaMin = omegaMin.Exp(two, big.NewInt(int64(bits-1)), nil).Add(omegaMin, one)
	omegaMax = omegaMax.Exp(two, big.NewInt(int64(bits)), nil).Sub(omegaMax, one)
	pi := GenCrtM(len)
	pi.Div(pi, two)
	tao, phi := new(big.Int), new(big.Int)
	tao.Mul(omegaMin, pi)
	phi.Mul(omegaMax, pi)

	start := time.Now()
	count := 0
	primesCount := 1
	for count < times {
	refresh:

		ck := new(big.Int)
		c, err := rand.Int(rand.Reader, phi)
		if err != nil {
			panic(fmt.Errorf("rand.Int failure in MJSpecialGenPrimes"))
		}
		ck.And(c, one)
		for !bool(ck.Cmp(one) == 0) {
			c, err := rand.Int(rand.Reader, phi)
			if err != nil {
				panic(fmt.Errorf("rand.Int failure in MJSpecialGenPrimes"))
			}
			ck.And(c, one)
		}

		ck.GCD(big.NewInt(0), big.NewInt(0), c, phi)

		for !bool(ck.Cmp(one) == 0) {
			c, err := rand.Int(rand.Reader, phi)
			if err != nil {
				panic(fmt.Errorf("rand.Int failure in MJSpecialGenPrimes"))
			}
			ck.GCD(big.NewInt(0), big.NewInt(0), c, phi)
		}

		q := new(big.Int)

	step:
		q.Add(c, tao)
		ck.And(q, big.NewInt(1))
		if !bool(ck.Cmp(one) == 0) {
			q.Add(q, pi)
			if q.BitLen() > bits {
				goto refresh
			}
		}

		for !q.ProbablyPrime(30) {
			primesCount++
			c.Mul(c, two)
			c.Mod(c, phi)
			goto step
		}
		count++
	}
	end := time.Since(start)
	fmt.Printf("MJ特例素数平均生成时间%fms\n", float64(end.Milliseconds())/float64(count))
	fmt.Printf("MJ特例素性测试次数%d\n", primesCount/count)
}

func MJImproveGenPrimes(times, bits, len int) {
	omegaMin, omegaMax := new(big.Int), new(big.Int)
	omegaMin = omegaMin.Exp(two, big.NewInt(int64(bits-1)), nil).Add(omegaMin, one)
	omegaMax = omegaMax.Exp(two, big.NewInt(int64(bits)), nil).Sub(omegaMax, one)

	omegaRange := new(big.Int)
	omegaRange.Sub(omegaMax, omegaMin)

	pi := GenCrtM(len)
	pi.Div(pi, two)

	taoTMax := new(big.Int)
	taoTMax.Sub(omegaMax, pi).Add(taoTMax, one)

	taoT := new(big.Int)
	taoTemp := new(big.Int)
	k := new(big.Int)
	ck := new(big.Int)

	count, primesCount := 0, 1
	start := time.Now()
	for count < times {
		ck.Set(zero)
		for !bool(ck.Cmp(one) == 0) {
			k, _ = rand.Int(rand.Reader, pi)
			ck.GCD(big.NewInt(0), big.NewInt(0), k, pi)
		}
	refresh:
		for taoT.Cmp(taoTMax) == 1 || taoT.Cmp(omegaMin) == -1 {
			v, err := rand.Int(rand.Reader, taoTemp.Div(taoTMax, pi))
			if err != nil {
				panic(fmt.Errorf("rand.Int failure in MJSpecialGenPrimes"))
			}
			taoT.Mul(v, pi)
		}

		p := new(big.Int)

		p.Add(k, taoT)

		if p.BitLen() > bits {
			goto refresh
		}

		for !p.ProbablyPrime(30) {
			primesCount++
			k.Mul(two, k).Mod(k, pi)
			goto refresh
		}
		count++
	}

	end := time.Since(start)
	fmt.Printf("改进的MJ素数平均生成时间%fms\n", float64(end.Milliseconds())/float64(count))
	fmt.Printf("改进的MJ素性测试次数%d\n", primesCount/count)
}

func ThresholdCRT512PrimesGen(threashold, times int) {
	n := big.NewInt(1)
	for i := 0; i < 73; i++ {
		n = n.Mul(n, big.NewInt(int64(primes1024[i])))
	}
	crt := make([]*big.Int, 73)
	crtPi := make([]*big.Int, 73)
	crtMi, crtMiIn := new(big.Int), new(big.Int)
	for i := 0; i < 73; i++ {
		crtPi[i] = big.NewInt(int64(primes1024[i]))
		crtMi.Div(n, crtPi[i])
		crtMiIn.ModInverse(crtMi, crtPi[i])
		crt[i] = new(big.Int)
		crt[i].Mul(crtMi, crtMiIn)
	}

	count, primesCount := 0, 1
	temp, p := new(big.Int), new(big.Int)
	start := time.Now()
	for count < times {
		temp.Set(zero)
		for temp.BitLen() == 0 || temp.BitLen() > 512 {
			temp.Set(zero)
			for i := 0; i < 73-threashold; i++ {
				crtRi, err := rand.Int(rand.Reader, crtPi[i])
				if err != nil {
					panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
				}
				for crtRi.BitLen() == 0 {
					crtRi, err = rand.Int(rand.Reader, crtPi[i])
					if err != nil {
						panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
					}
				}
				crtRi.Mul(crtRi, crt[i])
				temp.Add(temp, crtRi)
			}
			//fmt.Println(temp.BitLen())
		}
		p.Set(zero)
		for p.BitLen() != 512 || !p.ProbablyPrime(30) {
			p.Set(temp)
			for i := 73 - threashold; i < 73; i++ {
				crtRi, err := rand.Int(rand.Reader, crtPi[i])
				if err != nil {
					panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
				}
				for crtRi.BitLen() == 0 {
					crtRi, err = rand.Int(rand.Reader, crtPi[i])
					if err != nil {
						panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
					}
				}
				crtRi.Mul(crtRi, crt[i])
				p.Add(p, crtRi)
			}

			primesCount++
		}
		count++
	}
	end := time.Since(start)
	fmt.Printf("512门限CRT增量素数平均生成时间%fms\n", float64(end.Milliseconds())/float64(count))
	fmt.Printf("512素性测试次数%d\n", primesCount/count)
}

func ThresholdCRT1024PrimesGen(threashold, times int) {
	n := big.NewInt(1)
	for i := 0; i < 130; i++ {
		n = n.Mul(n, big.NewInt(int64(primes2048[i])))
	}
	crt := make([]*big.Int, 130)
	crtPi := make([]*big.Int, 130)
	crtMi, crtMiIn := new(big.Int), new(big.Int)
	for i := 0; i < 130; i++ {
		crtPi[i] = big.NewInt(int64(primes2048[i]))
		crtMi.Div(n, crtPi[i])
		crtMiIn.ModInverse(crtMi, crtPi[i])
		crt[i] = new(big.Int)
		crt[i].Mul(crtMi, crtMiIn)
	}

	count, primesCount := 0, 1
	temp, p := new(big.Int), new(big.Int)
	start := time.Now()
	for count < times {
		temp.Set(zero)
		for temp.BitLen() == 0 || temp.BitLen() > 1024 {
			temp.Set(zero)
			for i := 0; i < 130-threashold; i++ {
				crtRi, err := rand.Int(rand.Reader, crtPi[i])
				if err != nil {
					panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
				}
				for crtRi.BitLen() == 0 {
					crtRi, err = rand.Int(rand.Reader, crtPi[i])
					if err != nil {
						panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
					}
				}
				crtRi.Mul(crtRi, crt[i])
				temp.Add(temp, crtRi)
			}
		}
		p.Set(zero)
		for p.BitLen() != 1024 || !p.ProbablyPrime(30) {
			p.Set(temp)
			for i := 130 - threashold; i < 130; i++ {
				crtRi, err := rand.Int(rand.Reader, crtPi[i])
				if err != nil {
					panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
				}
				for crtRi.BitLen() == 0 {
					crtRi, err = rand.Int(rand.Reader, crtPi[i])
					if err != nil {
						panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
					}
				}
				crtRi.Mul(crtRi, crt[i])
				p.Add(p, crtRi)
			}
			primesCount++
		}
		count++
	}
	end := time.Since(start)
	fmt.Printf("1024门限CRT增量素数平均生成时间%fms\n", float64(end.Milliseconds())/float64(count))
	fmt.Printf("1024素性测试次数%d\n", primesCount/count)
}

func ThresholdCRT2048PrimesGen(threashold, times int) {
	n := big.NewInt(1)
	for i := 0; i < 231; i++ {
		n = n.Mul(n, big.NewInt(int64(primes4096[i])))
	}
	crt := make([]*big.Int, 231)
	crtPi := make([]*big.Int, 231)
	crtMi, crtMiIn := new(big.Int), new(big.Int)
	for i := 0; i < 231; i++ {
		crtPi[i] = big.NewInt(int64(primes4096[i]))
		crtMi.Div(n, crtPi[i])
		crtMiIn.ModInverse(crtMi, crtPi[i])
		crt[i] = new(big.Int)
		crt[i].Mul(crtMi, crtMiIn)
	}

	count, primesCount := 0, 1
	temp, p := new(big.Int), new(big.Int)
	start := time.Now()
	for count < times {
		temp.Set(zero)
		for temp.BitLen() == 0 || temp.BitLen() > 2048 {
			temp.Set(zero)
			for i := 0; i < 231-threashold; i++ {
				crtRi, err := rand.Int(rand.Reader, crtPi[i])
				if err != nil {
					panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
				}
				for crtRi.BitLen() == 0 {
					crtRi, err = rand.Int(rand.Reader, crtPi[i])
					if err != nil {
						panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
					}
				}
				crtRi.Mul(crtRi, crt[i])
				temp.Add(temp, crtRi)
			}
		}
		p.Set(zero)
		for p.BitLen() != 2048 || !p.ProbablyPrime(30) {
			p.Set(temp)
			for i := 231 - threashold; i < 231; i++ {
				crtRi, err := rand.Int(rand.Reader, crtPi[i])
				if err != nil {
					panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
				}
				for crtRi.BitLen() == 0 {
					crtRi, err = rand.Int(rand.Reader, crtPi[i])
					if err != nil {
						panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
					}
				}
				crtRi.Mul(crtRi, crt[i])
				p.Add(p, crtRi)
			}
			primesCount++
		}
		count++
	}
	end := time.Since(start)
	fmt.Printf("2048门限CRT增量素数平均生成时间%fms\n", float64(end.Milliseconds())/float64(count))
	fmt.Printf("2048素性测试次数%d\n", primesCount/count)
}

func ThresholdCRT512PrimesGenTest(threshold, times int) {
	n := big.NewInt(1)
	for i := 0; i < 73; i++ {
		n = n.Mul(n, big.NewInt(int64(primes1024[i])))
	}
	crt := make([]*big.Int, 73)
	crtPi := make([]*big.Int, 73)
	crtMi, crtMiIn := new(big.Int), new(big.Int)
	for i := 0; i < 73; i++ {
		crtPi[i] = big.NewInt(int64(primes1024[i]))
		crtMi.Div(n, crtPi[i])
		crtMiIn.ModInverse(crtMi, crtPi[i])
		crt[i] = new(big.Int)
		crt[i].Mul(crtMi, crtMiIn)
	}

	count, primesCount := 0, 1
	temp, p := new(big.Int), new(big.Int)
	start := time.Now()
	for count < times {
		temp.Set(zero)
		for temp.BitLen() == 0 || temp.BitLen() > 512 {
			temp.Set(zero)
			for i := 0; i < 73-threshold; i++ {
				crtRi, err := rand.Int(rand.Reader, crtPi[i])
				if err != nil {
					panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
				}
				for crtRi.BitLen() == 0 {
					crtRi, err = rand.Int(rand.Reader, crtPi[i])
					if err != nil {
						panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
					}
				}
				crtRi.Mul(crtRi, crt[i])
				temp.Add(temp, crtRi)
			}
		}
	re:
		for p.BitLen() > 512{
			p.Set(temp)
			for i := 73 - threshold; i < 73; i++ {
				crtRi, err := rand.Int(rand.Reader, crtPi[i])
				if err != nil {
					panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
				}
				for crtRi.BitLen() == 0 {
					crtRi, err = rand.Int(rand.Reader, crtPi[i])
					if err != nil {
						panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
					}
				}
				crtRi.Mul(crtRi, crt[i])
				p.Add(p, crtRi)
			}
		}

		for {
			if p.BitLen() > 512 {
				goto re
			}
			if p.BitLen() == 512 && p.ProbablyPrime(30) {
				primesCount++
				count++
				p.Set(zero)
				break
			}
			primesCount++
			p.Add(p, temp)
		}
	}
	end := time.Since(start)
	fmt.Printf("512特例门限CRT增量素数平均生成时间%fms\n", float64(end.Milliseconds())/float64(count))
	fmt.Printf("512特例素性测试次数%d\n", primesCount/count)
	fmt.Printf("512特例素性测试次数%d\n", primesCount)
}

func ThresholdCRT1024PrimesGenTest(threshold, times int) {
	n := big.NewInt(1)
	for i := 0; i < 130; i++ {
		n = n.Mul(n, big.NewInt(int64(primes2048[i])))
	}

	crt := make([]*big.Int, 130)
	crtPi := make([]*big.Int, 130)
	crtMi, crtMiIn := new(big.Int), new(big.Int)
	for i := 0; i < 130; i++ {
		crtPi[i] = big.NewInt(int64(primes2048[i]))
		crtMi.Div(n, crtPi[i])
		crtMiIn.ModInverse(crtMi, crtPi[i])
		crt[i] = new(big.Int)
		crt[i].Mul(crtMi, crtMiIn)
	}

	count, primesCount := 0, 1
	temp, p := new(big.Int), new(big.Int)
	start := time.Now()
	for count < times {
		temp.Set(zero)
		for temp.BitLen() == 0 || temp.BitLen() > 1024 {
			temp.Set(zero)
			for i := 0; i < 130-threshold; i++ {
				crtRi, err := rand.Int(rand.Reader, crtPi[i])
				if err != nil {
					panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
				}
				for crtRi.BitLen() == 0 {
					crtRi, err = rand.Int(rand.Reader, crtPi[i])
					if err != nil {
						panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
					}
				}
				crtRi.Mul(crtRi, crt[i])
				temp.Add(temp, crtRi)
			}
		}
	re:
		for p.BitLen() > 1024{
			p.Set(temp)
			for i := 130 - threshold; i < 130; i++ {
				crtRi, err := rand.Int(rand.Reader, crtPi[i])
				if err != nil {
					panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
				}
				for crtRi.BitLen() == 0 {
					crtRi, err = rand.Int(rand.Reader, crtPi[i])
					if err != nil {
						panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
					}
				}
				crtRi.Mul(crtRi, crt[i])
				p.Add(p, crtRi)
			}
		}

		for {
			if p.BitLen() > 1024 {
				goto re
			}
			if p.BitLen() == 1024 && p.ProbablyPrime(30) {
				primesCount++
				count++
				p.Set(zero)
				break
			}
			primesCount++
			p.Add(p, temp)
		}
	}
	end := time.Since(start)
	fmt.Printf("1024特例门限CRT增量素数平均生成时间%fms\n", float64(end.Milliseconds())/float64(count))
	fmt.Printf("1024特例素性测试次数%d\n", primesCount/count)
	fmt.Printf("1024特例素性测试次数%d\n", primesCount)
}

func ThresholdCRT2048PrimesGenTest(threshold, times int) {
	n := big.NewInt(1)
	for i := 0; i < 231; i++ {
		n = n.Mul(n, big.NewInt(int64(primes4096[i])))
	}
	crt := make([]*big.Int, 231)
	crtPi := make([]*big.Int, 231)
	crtMi, crtMiIn := new(big.Int), new(big.Int)
	for i := 0; i < 231; i++ {
		crtPi[i] = big.NewInt(int64(primes4096[i]))
		crtMi.Div(n, crtPi[i])
		crtMiIn.ModInverse(crtMi, crtPi[i])
		crt[i] = new(big.Int)
		crt[i].Mul(crtMi, crtMiIn)
	}

	count, primesCount := 0, 1
	temp, p := new(big.Int), new(big.Int)
	start := time.Now()
	for count < times {
		temp.Set(zero)
		for temp.BitLen() == 0 || temp.BitLen() > 2048 {
			temp.Set(zero)
			for i := 0; i < 231-threshold; i++ {
				crtRi, err := rand.Int(rand.Reader, crtPi[i])
				if err != nil {
					panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
				}
				for crtRi.BitLen() == 0 {
					crtRi, err = rand.Int(rand.Reader, crtPi[i])
					if err != nil {
						panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
					}
				}
				crtRi.Mul(crtRi, crt[i])
				temp.Add(temp, crtRi)
			}
		}
	re:
		for p.BitLen() > 2048{
			p.Set(temp)
			for i := 231 - threshold; i < 231; i++ {
				crtRi, err := rand.Int(rand.Reader, crtPi[i])
				if err != nil {
					panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
				}
				for crtRi.BitLen() == 0 {
					crtRi, err = rand.Int(rand.Reader, crtPi[i])
					if err != nil {
						panic(fmt.Errorf("rand.Int failure in ThresholdCRT512PrimesGen"))
					}
				}
				crtRi.Mul(crtRi, crt[i])
				p.Add(p, crtRi)
			}
		}

		for {
			if p.BitLen() == 2048 && p.ProbablyPrime(30) {
				primesCount++
				count++
				p.Set(zero)
				break
			}
			primesCount++
			p.Add(p, temp)
			if p.BitLen() > 2048 {
				goto re
			}
		}
	}
	end := time.Since(start)
	fmt.Printf("2048特例门限CRT增量素数平均生成时间%fms\n", float64(end.Milliseconds())/float64(count))
	fmt.Printf("2048特例素性测试次数%d\n", primesCount/count)
	fmt.Printf("2048特例素性测试次数%d\n", primesCount)
}

