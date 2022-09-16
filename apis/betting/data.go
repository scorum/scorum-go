package betting

import (
	"github.com/google/uuid"

	"github.com/scorum/scorum-go/types"
)

type Winner struct {
	Winner Better       `json:"winner"`
	Loser  Better       `json:"loser"`
	Market types.Market `json:"market"`
	Profit types.Asset  `json:"profit"`
	Income types.Asset  `json:"income"`
}

type Better struct {
	UUID        uuid.UUID     `json:"uuid"`
	AccountName string        `json:"name"`
	Wincase     types.Wincase `json:"wincase"`
}
