package sign

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	"github.com/pkg/errors"
)

func RefBlockNum(blockNumber uint32) uint16 {
	return uint16(blockNumber)
}

func RefBlockPrefix(blockID string) (uint32, error) {
	// Block GameID is hex-encoded.
	rawBlockID, err := hex.DecodeString(blockID)
	if err != nil {
		return 0, errors.Wrapf(err, "network_broadcast: failed to decode block GameID: %v", blockID)
	}

	// Raw prefix = raw block GameID [4:8].
	// Make sure we don't trigger a slice bounds out of range panic.
	if len(rawBlockID) < 8 {
		return 0, errors.Errorf("network_broadcast: invalid block GameID: %v", blockID)
	}
	rawPrefix := rawBlockID[4:8]

	// Decode the prefix.
	var prefix uint32
	if err := binary.Read(bytes.NewReader(rawPrefix), binary.LittleEndian, &prefix); err != nil {
		return 0, errors.Wrapf(err, "network_broadcast: failed to read block prefix: %v", rawPrefix)
	}

	// Done, return the prefix.
	return prefix, nil
}
