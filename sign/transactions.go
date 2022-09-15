package sign

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

func RefBlockNum(blockNumber uint32) uint16 {
	return uint16(blockNumber)
}

func RefBlockPrefix(blockID string) (uint32, error) {
	// Block ID is hex-encoded.
	rawBlockID, err := hex.DecodeString(blockID)
	if err != nil {
		return 0, fmt.Errorf("failed to decode block ID: %w", err)
	}

	// Raw prefix = raw block ID [4:8].
	// Make sure we don't trigger a slice bounds out of range panic.
	if len(rawBlockID) < 8 {
		return 0, errors.New("invalid blockID")
	}
	rawPrefix := rawBlockID[4:8]

	// Decode the prefix.
	var prefix uint32
	if err := binary.Read(bytes.NewReader(rawPrefix), binary.LittleEndian, &prefix); err != nil {
		return 0, fmt.Errorf("failed to read block prefix: %w", err)
	}

	// Done, return the prefix.
	return prefix, nil
}
