package sign

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/scorum/scorum-go/types"
	"github.com/stretchr/testify/require"
)

func TestTransaction_Digest(t *testing.T) {
	var tx *types.Transaction
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

	var initChain = &Chain{
		ID: "0000000000000000000000000000000000000000000000000000000000000000",
	}

	expected := "582176b1daf89984bc8b4fdcb24ff1433d1eb114a8c4bf20fb22ad580d035889"
	stx := NewSignedTransaction(tx)
	digest, err := stx.Digest(initChain)
	require.NoError(t, err)
	require.Equal(t, expected, hex.EncodeToString(digest))
}

func TestTransaction_Verify(t *testing.T) {
	var tx *types.Transaction
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

	var initChain = &Chain{
		ID: "0000000000000000000000000000000000000000000000000000000000000000",
	}

	expectedDigest := "582176b1daf89984bc8b4fdcb24ff1433d1eb114a8c4bf20fb22ad580d035889"
	stx := NewSignedTransaction(tx)
	digest, err := stx.Digest(initChain)
	require.NoError(t, err)
	require.Equal(t, expectedDigest, hex.EncodeToString(digest))

	wif := "5HzA7Eju6BQcunCHNy5ywGG2YriZMQaEcoPEbQFpSRk9NZEo6Fv"
	require.NoError(t, stx.Sign([]string{wif}, initChain))

	btcecWif, err := btcutil.DecodeWIF(wif)
	require.NoError(t, err)

	res, err := stx.Verify(initChain, [][]byte{btcecWif.PrivKey.PubKey().SerializeUncompressed()})
	require.NoError(t, err)
	require.True(t, res)
}
