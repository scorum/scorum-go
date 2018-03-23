package sign

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/scorum/scorum-go/encoding/wif"
	"github.com/scorum/scorum-go/types"
	"github.com/stretchr/testify/require"
)

var initChain = &Chain{
	ID: "0000000000000000000000000000000000000000000000000000000000000000",
}

var tx *types.Transaction

func init() {
	// Prepare the transaction.
	expiration := time.Date(2016, 8, 8, 12, 24, 17, 0, time.UTC)
	tx = &types.Transaction{
		RefBlockNum:    36029,
		RefBlockPrefix: 1164960351,
		Expiration:     &types.Time{&expiration},
	}
	tx.PushOperation(&types.VoteOperation{
		Voter:    "xeroc",
		Author:   "xeroc",
		Permlink: "piston",
		Weight:   10000,
	})
}

var wifs = []string{
	"5JLw5dgQAx6rhZEgNN5C2ds1V47RweGshynFSWFbaMohsYsBvE8",
}

var privateKeys = make([][]byte, 0, len(wifs))

func init() {
	for _, v := range wifs {
		privKey, err := wif.Decode(v)
		if err != nil {
			panic(err)
		}
		privateKeys = append(privateKeys, privKey)
	}
}

var publicKeys = make([][]byte, 0, len(wifs))

func init() {
	for _, v := range wifs {
		pubKey, err := wif.GetPublicKey(v)
		if err != nil {
			panic(err)
		}
		publicKeys = append(publicKeys, pubKey)
	}
}

func TestTransaction_Digest(t *testing.T) {
	expected := "582176b1daf89984bc8b4fdcb24ff1433d1eb114a8c4bf20fb22ad580d035889"

	stx := NewSignedTransaction(tx)

	digest, err := stx.Digest(initChain)
	require.NoError(t, err)
	require.Equal(t, expected, hex.EncodeToString(digest))
}

func TestTransaction_SignAndVerify(t *testing.T) {
	tx.Signatures = nil
	defer func() {
		tx.Signatures = nil
	}()

	stx := NewSignedTransaction(tx)
	require.NoError(t, stx.Sign(privateKeys, initChain))
	require.Len(t, tx.Signatures, 1, "expected signatures not appended to the transaction")

	t.Log("signature", stx.Signatures)

	ok, err := stx.Verify(publicKeys, initChain)
	require.NoError(t, err)
	require.True(t, ok, "verification failed")
}
