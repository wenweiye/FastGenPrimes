package main

import "Go-EC/common"

func main(){
	//common.TraPrimeGen(100000,512)
	//common.TraPrimeGen(100000,1024)
	//common.TraPrimeGen(100000,2048)
	common.TraPrimeGenWithIncreasing(100000,512)
	//common.TraPrimeGenWithIncreasing(100000,1024)
	//common.TraPrimeGenWithIncreasing(100000,2048)
	//common.TraPrimeGenWithImProIncreasing(100000,512,70)
	//common.TraPrimeGenWithImProIncreasing(100000,1024,70)
	//common.TraPrimeGenWithImProIncreasing(100000,2048,70)
	//common.MJSpecialGenPrimes(100000,512,10)
	//common.MJSpecialGenPrimes(100000,1024,130)
	//common.MJSpecialGenPrimes(100000,2048,130)
	//common.ThresholdCRT512PrimesGen(3,100000)
	//common.ThresholdCRT1024PrimesGen(3,100000)
	//common.ThresholdCRT2048PrimesGen(3,100000)
	common.ThresholdCRT512PrimesGenTest(3,100)
	common.ThresholdCRT1024PrimesGenTest(3,100)
	common.ThresholdCRT2048PrimesGenTest(3,100)
	//common.MJImproveGenPrimes(100000,512,10)
	//common.MJImproveGenPrimes(100000,1024,130)
	//common.MJImproveGenPrimes(100000,2048,130)
}
