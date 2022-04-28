package sign

import (
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/scorum/scorum-go/key"
	"github.com/scorum/scorum-go/types"
)

var (
	zeroChainID, _ = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
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

	expected := "582176b1daf89984bc8b4fdcb24ff1433d1eb114a8c4bf20fb22ad580d035889"
	stx := NewSignedTransaction(tx)
	digest, err := stx.Digest(zeroChainID)
	require.NoError(t, err)
	require.Equal(t, expected, hex.EncodeToString(digest))
}

func TestTransaction_Verify(t *testing.T) {
	txStr := `
    {
        "ref_block_num": 57753,
        "ref_block_prefix": 1882698764,
        "expiration": "2021-12-01T17:26:00",
        "operations": [
            [
                "transfer",
                {
                    "from": "azucena",
                    "to": "leonarda",
                    "amount": "0.000009000 SCR",
                    "memo": "{\"bet_id\":\"8b022219-3825-413e-a6ae-1cc3154bdb7f\",\"game_id\":\"17ff54ad-8472-4e4f-9e96-2b6ccbfaef33\"}"
                }
            ]
        ],
        "extensions": [],
        "signatures": [
            "204da1d0e0c5cb39978d1ad4fe4cd4446718ad0e10a19d662ac7bc9e0bba708f2c2116279fdc512d03fe2c97cdfc00804a4b958a8078f9b3b28e4f4b22f798f032"
        ],
        "transaction_id": "825d2733272648f95c21ae97be9f8c3d8edaf8f9",
        "block_num": 31449514,
        "transaction_num": 0
    }`

	var tx types.Transaction
	err := json.Unmarshal([]byte(txStr), &tx)
	stx := NewSignedTransaction(&tx)

	pubKey, err := key.NewPublicKey("SCR7cTf2Dx9rxffs6E2z2pdn5cLMneo3AAFSsF9g4SaVviCYdfQ63")
	require.NoError(t, err)

	require.NoError(t, stx.Verify(TestNetChainID, pubKey))
}
