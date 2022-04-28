package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/scorum/scorum-go/encoding/transaction"

	"github.com/stretchr/testify/require"
)

func TestPublicKey_MarshalTransaction(t *testing.T) {
	var (
		b       bytes.Buffer
		key     PublicKey = "SCR5jPZF7PMgTpLqkdfpMu8kXea8Gio6E646aYpTgcjr9qMLrAgnL"
		encoder           = transaction.NewEncoder(&b)
	)

	require.NoError(t, key.MarshalTransaction(encoder))
	require.Equal(t, "026f0896f24d94252c351715bfe6052bbf9ea820e805bd47c2496c626d3467da5d", hex.EncodeToString(b.Bytes()))
}
