package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/scorum/scorum-go/encoding/transaction"

	"github.com/stretchr/testify/require"
)

func TestWIFFromBytes(t *testing.T) {
	brainKey := "gadger varna mormon leet archway chewer yank outbuzz mailbag region douc upshut basenji blankly moineau treen"
	hash := sha256.Sum256([]byte(brainKey))

	wif, err := WIFFromBytes(hash[:])
	require.NoError(t, err)
	require.Equal(t, "5KAZ2g22NCiM815YqjLGFYDqwBF4f2CYA9cFhhg9qg9Zf4xtAfh", wif.String())
}

func TestWIFFromBrainKey(t *testing.T) {
	brainKey := "gadger varna mormon leet archway chewer yank outbuzz mailbag region douc upshut basenji blankly moineau treen"
	wif, err := WIFFromBrainKey(brainKey)
	require.NoError(t, err)
	require.Equal(t, "5KAZ2g22NCiM815YqjLGFYDqwBF4f2CYA9cFhhg9qg9Zf4xtAfh", wif.String())
}

func TestPublicKey_MarshalTransaction(t *testing.T) {
	t.Run("empty key", func(t *testing.T) {
		var key PublicKey
		var b bytes.Buffer
		encoder := transaction.NewEncoder(&b)
		require.EqualError(t, key.MarshalTransaction(encoder), "wrong prefix")
		require.Equal(t, "", hex.EncodeToString(b.Bytes()))
	})

	t.Run("invalid key", func(t *testing.T) {
		key := PublicKey("123")
		var b bytes.Buffer
		encoder := transaction.NewEncoder(&b)
		require.EqualError(t, key.MarshalTransaction(encoder), "wrong prefix")
		require.Equal(t, "", hex.EncodeToString(b.Bytes()))
	})

	t.Run("valid key", func(t *testing.T) {
		key := PublicKey("SCR5jPZF7PMgTpLqkdfpMu8kXea8Gio6E646aYpTgcjr9qMLrAgnL")
		var b bytes.Buffer
		encoder := transaction.NewEncoder(&b)
		require.NoError(t, key.MarshalTransaction(encoder))
		require.Equal(t, "026f0896f24d94252c351715bfe6052bbf9ea820e805bd47c2496c626d3467da5d", hex.EncodeToString(b.Bytes()))
	})
}
