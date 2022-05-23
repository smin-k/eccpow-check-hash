package main // import "github.com/smin-k/checkhash"

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/smin-k/checkhash/eccpow"

	"github.com/Onther-Tech/go-ethereum/common"
	"github.com/Onther-Tech/go-ethereum/core/types"
	//C:\Users\infonet\go\pkg\mod\github.com\!onther-!tech\go-ethereum@v1.8.23
)

type Header struct {
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	UncleHash   common.Hash    `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    common.Address `json:"miner"            gencodec:"required"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       types.Bloom    `json:"logsBloom"        gencodec:"required"`
	Difficulty  big.Int       `json:"difficulty"       gencodec:"required"`
	Number      big.Int       `json:"number"           gencodec:"required"`
	GasLimit    uint64         `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64         `json:"gasUsed"          gencodec:"required"`
	Time        uint64        `json:"timestamp"        gencodec:"required"`
	Extra       []byte         `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash    `json:"mixHash"          gencodec:"required"`
	Hash        common.Hash    `json:"hash"`
	Nonce       types.BlockNonce     `json:"nonce"            gencodec:"required"`
}

type verifyParameters struct {
	n          uint64
	m          uint64
	wc         uint64
	wr         uint64
	seed       uint64
	outputWord []uint64
}


func main() {
	ecc := eccpow.ECC{}
	header1, hashinblock:= get_header("json/block6209.json")
	hash1 := ecc.SealHash(&header1).Bytes()

	fmt.Printf("\nheader:\n%v\n", header1)
	fmt.Printf("\nblockNumber:%v\n", header1.Number)
	fmt.Printf("\nhashinblock:%v\n", hashinblock)
	fmt.Printf("\nhash:%v\n", hash1)

	bool, hashVectorOfVerification, outputWordOfVerification, seed := eccpow.VerifyOptimizedDecoding(&header1, hash1)
	
	fmt.Printf("bool: %v_____\n", bool)
	fmt.Printf("%v_____\n", hashVectorOfVerification)
	fmt.Printf("%v_____\n", outputWordOfVerification)
	fmt.Printf("seed: %v_____\n\n", seed)	

	fmt.Printf("seedTodigest: %v\n", common.BytesToHash(seed))
	fmt.Printf("digest: %v\n", header1.MixDigest)

}

func get_header(filepath string) (types.Header, common.Hash){
	//block 6209 in ETH-ECC1
	jsonFile, err := os.Open(filepath)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	
	if err != nil {
		fmt.Println(err)
	}
	
	var result Header
	var result2 types.Header
	var result3 map[string]interface{}

	json.Unmarshal(byteValue, &result)
	json.Unmarshal(byteValue, &result3)

	hexstring, _ := hex.DecodeString(result3["extraData"].(string)[2:])

	/*
	fmt.Printf("\nresult3extraData:%v\n", result3["extraData"])
	fmt.Printf("\nhexstring:%v\n", hexstring)
	*/
	
	result2.ParentHash = result.ParentHash
	result2.UncleHash = result.UncleHash
	result2.Coinbase = result.Coinbase
	result2.Root = result.Root
	result2.TxHash = result.TxHash
	result2.ReceiptHash = result.ReceiptHash
	result2.Bloom = result.Bloom
	result2.Difficulty = &result.Difficulty
	result2.Number = &result.Number
	result2.GasLimit = result.GasLimit
	result2.GasUsed = result.GasUsed
	result2.Time = result.Time
	result2.Extra = hexstring

	result2.MixDigest = result.MixDigest
	result2.Nonce = result.Nonce
	hash := result.Hash

	//fmt.Printf("_____%v_____\n\n\n", result2)

	return result2, hash
}

