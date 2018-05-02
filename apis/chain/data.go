package chain

import "github.com/scorum/scorum-go/types"

type ChainProperties struct {
	AverageBlockSize            uint32               `json:"average_block_size"`
	CurrentReserveRatio         uint32               `json:"current_reserve_ratio"`
	MaxVirtualBandwidth         string               `json:"max_virtual_bandwidth"`
	ChainID                     string               `json:"chain_id"`
	HeadBlockID                 string               `json:"head_block_id"`
	HeadBlockNumber             uint32               `json:"head_block_number"`
	LastIrreversibleBlockNumber uint32               `json:"last_irreversible_block_number"`
	CurrentAslot                uint32               `json:"current_aslot"`
	Time                        types.Time           `json:"time"`
	CurrentWitness              string               `json:"current_witness"`
	MedianChainProperies        MediaChainProperties `json:"median_chain_props"`
	MajorityVersion             string               `json:"majority_version"`
	HFVersion                   string               `json:"hf_version"`
}

type MediaChainProperties struct {
	AccountCreationFee types.Asset `json:"account_creation_fee"`
	MaximumBlockSize   uint32      `json:"maximum_block_size"`
}
