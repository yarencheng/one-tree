package block

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"yarencheng/one-tree/go-src/protobuf"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

func GenFirstBlock(privateKey *ecdsa.PrivateKey) protobuf.Block {

	payload := protobuf.Payload{
		LastUpdated: ptypes.TimestampNow(),
		Height:      0,
		Childs:      0,
	}

	payloadBytes, err := proto.Marshal(&payload)
	if err != nil {
		log.Fatalln("Failed to encode block payload: ", err)
	}

	hash := crypto.Keccak256Hash(payloadBytes)

	_, sig, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		log.Fatalln("Failed to sign: ", err)
	}

	block := protobuf.Block{
		PublicKey: &protobuf.PublicKey{
			Bytes: elliptic.Marshal(privateKey.PublicKey.Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y),
		},
		Signature: sig.Bytes(),
		Payload:   &payload,
	}

	return block
}
