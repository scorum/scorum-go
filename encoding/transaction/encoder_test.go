package transaction

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoder_EncodeMoney_ManyZero(t *testing.T) {
	buffer := new(bytes.Buffer)

	encoder := NewEncoder(buffer)
	err := encoder.EncodeMoney("00000000000000000000099.0000000000000000000000000 SCR")
	require.Error(t, err)
}

func TestEncoder_EncodeMoney_ValueOverflow(t *testing.T) {
	buffer := new(bytes.Buffer)
	encoder := NewEncoder(buffer)

	err := encoder.EncodeMoney("11111111111111111111111111111111111111 SCR")
	require.Error(t, err)
}
