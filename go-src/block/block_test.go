package block

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"testing"

	"encoding/hex"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestGenFirstBlock(t *testing.T) {

	// arrange

	privateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			X:     &big.Int{},
			Y:     &big.Int{},
			Curve: elliptic.P256(),
		},
		D: &big.Int{},
	}

	_, ok := privateKey.X.SetString("115374667210361068292087618343890579933223980850028273742606087063871943245543", 10)
	require.True(t, ok)
	_, ok = privateKey.Y.SetString("47917581441659907363930539417776011972510548999569359844922156418221962671735", 10)
	require.True(t, ok)
	_, ok = privateKey.D.SetString("9157388670700653708049816582433035038265863790928689767690366181005819786793", 10)
	require.True(t, ok)

	_, err := hex.DecodeString("04ff13bf5c8079a7dae3f7cb401e2f4ac3a6c01b6e67b69dab7dc23a7d6b6a7ee769f063d6183e7f909d7c1ebaa42e47c1960dc791250a165ffb9e73f25123b677")
	require.NoError(t, err)

	_, err = hex.DecodeString("195224f931f17ae9faa56f44508dbd9b3c015796a240d69564990f64f5750be9")
	require.NoError(t, err)

	// action

	block := GenFirstBlock(privateKey)
	_, err = proto.Marshal(block.Payload)
	require.NoError(t, err)

	// assert

	// assert.Equal(t, protobuf.Block{
	// 	PublicKey: &protobuf.PublicKey{
	// 		Bytes: pubKeyBytes,
	// 	},
	// 	Signature: signature,
	// }, block)
}
