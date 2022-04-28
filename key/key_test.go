package key

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/scorum/scorum-go/encoding/transaction"
	"github.com/stretchr/testify/require"
)

const (
	privateKeyStr = "5JWHY5DxTF6qN5grTtChDCYBmWHfY9zaSsw4CxEKN5eZpH9iBma"
	privateKeyHex = "5ad2b8df2c255d4a2996ee7d065e013e1bbb35c075ee6e5208aca44adc9a9d4c"
	publicKeyStr  = "SCR7jNh5ejQoqHqWcGWFJ1v4F5CzsG3EiBuz1VooCng1cH5QpJD27"
	publicKeyHex  = "0376645292a6ab11c53075bee9905afbc7168d7dec4260c2e9942abd92644de8ed"
	hashHex       = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	sigHex        = "207db2f1826ff4639599fd00f1f34df7fea05f60fd629f45f290b0fe21e12be8796dd42888b494ff163a6965385f2981edee609cebe6bf76846780ecc8870528ec"
)

func TestPublicKey(t *testing.T) {
	t.Run("empty key", func(t *testing.T) {
		var key PublicKey
		var b bytes.Buffer
		encoder := transaction.NewEncoder(&b)
		require.EqualError(t, key.MarshalTransaction(encoder), "public key is nil")
		require.Equal(t, "", hex.EncodeToString(b.Bytes()))
	})

	t.Run("wrong key len", func(t *testing.T) {
		_, err := NewPublicKey("")
		require.ErrorIs(t, err, ErrWrongKeyLen)
	})

	t.Run("wrong prefix", func(t *testing.T) {
		_, err := NewPublicKey("SSS11111111111111111111111111111111111111111111111111")
		require.ErrorIs(t, err, ErrWrongPrefix)
	})

	t.Run("wrong check sum", func(t *testing.T) {
		_, err := NewPublicKey("SCR11111111111111111111111111111111111111111111111111")
		require.ErrorIs(t, err, ErrWrongChecksum)
	})

	t.Run("public key from string", func(t *testing.T) {
		key, err := NewPublicKey(publicKeyStr)
		require.NoError(t, err)
		require.Equal(t, publicKeyStr, key.String())
		require.Equal(t, publicKeyHex, hex.EncodeToString(key.Serialize()))
	})

	t.Run("public key from bytes", func(t *testing.T) {
		raw, err := hex.DecodeString(publicKeyHex)
		require.NoError(t, err)

		key, err := PublicKeyFromBytes(raw)
		require.NoError(t, err)
		require.Equal(t, publicKeyStr, key.String())
		require.Equal(t, publicKeyHex, hex.EncodeToString(key.Serialize()))
	})
}

func TestNewPrivateKey(t *testing.T) {
	pk, err := NewPrivateKey()
	require.NoError(t, err)

	pk2, err := PrivateKeyFromBytes(pk.Serialize())
	require.NoError(t, err)

	pk3, err := PrivateKeyFromString(pk.String())
	require.NoError(t, err)

	require.Equal(t, pk.String(), pk2.String())
	require.Equal(t, pk.String(), pk3.String())
}

func TestPrivateKeyFromString(t *testing.T) {
	pk, err := PrivateKeyFromString(privateKeyStr)
	require.NoError(t, err)
	require.Equal(t, privateKeyHex, hex.EncodeToString(pk.Serialize()))
	require.Equal(t, publicKeyStr, pk.PublicKey().String())
}

func TestPrivateKeyFromBytes(t *testing.T) {
	raw, err := hex.DecodeString(privateKeyHex)
	require.NoError(t, err)
	pk, err := PrivateKeyFromBytes(raw)
	require.NoError(t, err)
	require.Equal(t, publicKeyStr, pk.PublicKey().String())
}

func TestSignAndValidate(t *testing.T) {
	pk, err := PrivateKeyFromString(privateKeyStr)
	require.NoError(t, err)

	hash, err := hex.DecodeString(hashHex)
	require.NoError(t, err)

	sig := pk.Sign(hash)
	require.Equal(t, sigHex, hex.EncodeToString(sig))

	err = pk.PublicKey().Verify(hash, sig)
	require.NoError(t, err)

	t.Run("validate with wrong key key mismatch", func(t *testing.T) {
		wrogKey, err := NewPublicKey("SCR5jPZF7PMgTpLqkdfpMu8kXea8Gio6E646aYpTgcjr9qMLrAgnL")
		require.NoError(t, err)
		err = wrogKey.Verify(hash, sig)
		require.ErrorIs(t, err, ErrKeyMismatch)
	})
}
