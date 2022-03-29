package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/scorum/scorum-go/encoding/transaction"
	"github.com/stretchr/testify/require"
)

func TestKeyAuthorityMap_MarshalTransaction(t *testing.T) {
	v := NewKeyAuthorityMap(KeyAuthority{Key: "SCR7zPNg5nAsJjP9gvMfQ4UnAwDwf91WPYC8KFzobtMuQ52ns1D6T", Weight: 1})

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)
	require.NoError(t, v.MarshalTransaction(encoder))
	require.Equal(t, "0103987a5a967458c114c15091198c06a822f54b494ea486204551a53f85effa31420100", hex.EncodeToString(b.Bytes()))
}

func TestKeyAuthorityMap_MarshalJSON(t *testing.T) {
	v := NewKeyAuthorityMap(KeyAuthority{Key: "SCR7zPNg5nAsJjP9gvMfQ4UnAwDwf91WPYC8KFzobtMuQ52ns1D6T", Weight: 1})

	d, err := json.Marshal(v)
	require.NoError(t, err)
	require.Equal(t, "[[\"SCR7zPNg5nAsJjP9gvMfQ4UnAwDwf91WPYC8KFzobtMuQ52ns1D6T\",1]]", string(d))
}

func TestAccountAuthorityMap(t *testing.T) {
	v := NewAccountAuthorityMap(AccountAuthority{AccountName: "alice", Weight: 1})

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)
	require.NoError(t, v.MarshalTransaction(encoder))
	require.Equal(t, "0105616c6963650100", hex.EncodeToString(b.Bytes()))
}
