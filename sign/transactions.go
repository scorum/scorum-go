package sign

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func ParseBlockID(blockID string) (uint16, uint32, error) {
	// Block ID is hex-encoded.
	rawBlockID, err := hex.DecodeString(blockID)
	if err != nil {
		return 0, 0, fmt.Errorf("network_broadcast: failed to decode block ID: %v: %w", blockID, err)
	}

	// Raw prefix = raw block ID [4:8].
	// Make sure we don't trigger a slice bounds out of range panic.
	if len(rawBlockID) < 8 {
		return 0, 0, fmt.Errorf("network_broadcast: invalid block ID: %v", blockID)
	}

	rawNum := rawBlockID[:4]
	rawPrefix := rawBlockID[4:8]

	// Decode the prefix.
	var prefix uint32
	if err := binary.Read(bytes.NewReader(rawPrefix), binary.LittleEndian, &prefix); err != nil {
		return 0, 0, fmt.Errorf("network_broadcast: failed to read block prefix: %v: %w", rawPrefix, err)
	}

	var num uint32
	if err := binary.Read(bytes.NewReader(rawNum), binary.BigEndian, &num); err != nil {
		return 0, 0, fmt.Errorf("network_broadcast: failed to read block number: %v: %w", rawNum, err)
	}

	return uint16(num), prefix, nil
}
