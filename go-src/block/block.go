package block

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"yarencheng/one-tree/go-src/pb"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

func GenFirstBlock(privateKey *ecdsa.PrivateKey) pb.Block {

	payload := pb.Payload{
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

	block := pb.Block{
		PublicKey: &pb.PublicKey{
			Bytes: elliptic.Marshal(privateKey.PublicKey.Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y),
		},
		Signature: sig.Bytes(),
		Payload:   &payload,
	}

	return block
}
