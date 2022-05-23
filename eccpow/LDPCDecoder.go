package eccpow

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/Onther-Tech/go-ethereum/core/types"
	"github.com/Onther-Tech/go-ethereum/crypto"
	"github.com/Onther-Tech/go-ethereum/log"
)

//OptimizedDecoding return hashVector, outputWord, LRrtl
func OptimizedDecoding(parameters Parameters, hashVector []int, H, rowInCol, colInRow [][]int) ([]int, []int, [][]float64) {
	outputWord := make([]int, parameters.n)
	LRqtl := make([][]float64, parameters.n)
	LRrtl := make([][]float64, parameters.n)
	LRft := make([]float64, parameters.n)

	for i := 0; i < parameters.n; i++ {
		LRqtl[i] = make([]float64, parameters.m)
		LRrtl[i] = make([]float64, parameters.m)
		LRft[i] = math.Log((1-crossErr)/crossErr) * float64((hashVector[i]*2 - 1))
	}
	LRpt := make([]float64, parameters.n)

	for ind := 1; ind <= maxIter; ind++ {
		for t := 0; t < parameters.n; t++ {
			temp3 := 0.0

			for mp := 0; mp < parameters.wc; mp++ {
				temp3 = infinityTest(temp3 + LRrtl[t][rowInCol[mp][t]])
			}
			for m := 0; m < parameters.wc; m++ {
				temp4 := temp3
				temp4 = infinityTest(temp4 - LRrtl[t][rowInCol[m][t]])
				LRqtl[t][rowInCol[m][t]] = infinityTest(LRft[t] + temp4)
			}
		}

		for k := 0; k < parameters.wr; k++ {
			for l := 0; l < parameters.wr; l++ {
				temp3 := 0.0
				sign := 1.0
				tempSign := 0.0
				for m := 0; m < parameters.wr; m++ {
					if m != l {
						temp3 = temp3 + funcF(math.Abs(LRqtl[colInRow[m][k]][k]))
						if LRqtl[colInRow[m][k]][k] > 0.0 {
							tempSign = 1.0
						} else {
							tempSign = -1.0
						}
						sign = sign * tempSign
					}
				}
				magnitude := funcF(temp3)
				LRrtl[colInRow[l][k]][k] = infinityTest(sign * magnitude)
			}
		}

		for t := 0; t < parameters.n; t++ {
			LRpt[t] = infinityTest(LRft[t])
			for k := 0; k < parameters.wc; k++ {
				LRpt[t] += LRrtl[t][rowInCol[k][t]]
				LRpt[t] = infinityTest(LRpt[t])
			}

			if LRpt[t] >= 0 {
				outputWord[t] = 1
			} else {
				outputWord[t] = 0
			}
		}
	}

	return hashVector, outputWord, LRrtl
}

//VerifyOptimizedDecoding return bool, hashVector, outputword, digest which are used for validation
func VerifyOptimizedDecoding(header *types.Header, hash []byte) (bool, []int, []int, []byte) {

	parameters, _ := setParameters(header)
	H, _, _ := generateH(parameters)
	fmt.Printf("Parameters:_%v\n", parameters)
	fmt.Printf("Parameters.seed:_%v\n", parameters.seed)
	//fmt.Printf("H:_\n%v\n", H)
	//fmt.Printf("hSeed:_\n%v\n", hSeed)
	//fmt.Printf("colOrder:_\n%v\n", colOrder)
	colInRow, rowInCol := generateQ(parameters, H)

	seed := make([]byte, 40)
	copy(seed, hash)
	//fmt.Printf("\nNonce:_%v_____\n", header.Nonce)
	
	binary.LittleEndian.PutUint64(seed[32:], header.Nonce.Uint64())
	//fmt.Printf("\nseed:_%v_____\n", seed)

	seed = crypto.Keccak512(seed)

	//fmt.Printf("\nKeccak_seed:_%v_____\n", seed)

	hashVector := generateHv(parameters, seed)
	hashVectorOfVerification, outputWordOfVerification, _ := OptimizedDecoding(parameters, hashVector, H, rowInCol, colInRow)

	if MakeDecision(header, colInRow, outputWordOfVerification) {
		log.Warn("success: Verify")
		return true, hashVectorOfVerification, outputWordOfVerification, seed
	}

	return false, hashVectorOfVerification, outputWordOfVerification, seed
}
