package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsset_UnmarshalJSON(t *testing.T) {
	var asset Asset
	err := asset.UnmarshalText([]byte("1.3003 SCR"))
	require.NoError(t, err)

	require.Equal(t, "1.300300000 SCR", asset.String())
}

func TestAsset_MarshalJSON(t *testing.T) {
	asset := AssetFromFloat(123.56)
	bytes, err := asset.MarshalText()
	require.NoError(t, err)
	require.Equal(t, `123.560000000 SCR`, string(bytes))
}

func TestAsset_String(t *testing.T) {
	asset := AssetFromFloat(123.56)
	require.Equal(t, "123.560000000 SCR", asset.String())
}

func TestAsset_AssertFromString(t *testing.T) {
	asset, err := AssetFromString("123.56 SCR")
	require.NoError(t, err)
	require.Equal(t, "123.560000000 SCR", asset.String())

	asset, err = AssetFromString("123.56 SCP")
	require.Error(t, err)

	asset, err = AssetFromString("123.56 SCR PTR")
	require.Error(t, err)
}
