package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"yarencheng/one-tree/go-src/block"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	fmt.Println("Hello blockchain")

	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	b := block.SignFirstBlock(privateKey)

	fmt.Println("b = %#v", b)

	pb := fmt.Sprintf("%#v", b.Payload)
	hash := sha256.Sum256([]byte(pb))

	valid := crypto.VerifySignature(crypto.FromECDSAPub(&privateKey.PublicKey), hash[:], b.Signature)

	fmt.Println("valid = %#v", valid)
}
