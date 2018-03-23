package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsset_UnmarshalJSON(t *testing.T) {
	var asset Asset
	err := asset.UnmarshalJSON([]byte("\"1.3003 SCR\""))
	require.NoError(t, err)

	require.Equal(t, "1.300300000 SCR", asset.String())
}

func TestAsset_MarshalJSON(t *testing.T) {
	asset := AssertFromFloat(123.56)
	bytes, err := asset.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"123.560000000 SCR"`, string(bytes))
}

func TestAsset_String(t *testing.T) {
	asset := AssertFromFloat(123.56)
	require.Equal(t, "123.560000000 SCR", asset.String())
}
