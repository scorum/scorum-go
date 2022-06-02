package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/scorum/scorum-go/encoding/transaction"
	"github.com/stretchr/testify/require"
)

func TestCreateGameOperation_SerializationWithoutMarkets(t *testing.T) {
	time, err := time.Parse(Layout, `"2018-08-03T10:12:43"`)
	require.NoError(t, err)

	uuid := uuid.UUID{}
	require.NoError(t, uuid.UnmarshalText([]byte("e629f9aa-6b2c-46aa-8fa8-36770e7a7a5f")))

	op := CreateGameOperation{
		UUID:                uuid,
		Moderator:           "admin",
		JsonMetadata:        "{}",
		StartTime:           Time{&time},
		GameType:            SoccerGameType,
		AutoResolveDelaySec: 33,
	}

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)
	require.NoError(t, op.MarshalTransaction(encoder))

	require.EqualValues(t, "23e629f9aa6b2c46aa8fa836770e7a7a5f0561646d696e027b7d9b2a645b210000000000", hex.EncodeToString(b.Bytes()))
}

func TestCreateGameOperation_SerializationWithTotalMarket(t *testing.T) {
	time, err := time.Parse(Layout, `"2018-08-03T10:12:43"`)
	require.NoError(t, err)

	uuid := uuid.UUID{}
	require.NoError(t, uuid.UnmarshalText([]byte("e629f9aa-6b2c-46aa-8fa8-36770e7a7a5f")))

	op := CreateGameOperation{
		UUID:                uuid,
		Moderator:           "admin",
		JsonMetadata:        "{}",
		StartTime:           Time{&time},
		AutoResolveDelaySec: 33,
		Markets: []Market{Market{&OverUnderMarket{
			ID:        MarketTotal,
			Threshold: 1000,
		}}},
	}

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)
	require.NoError(t, op.MarshalTransaction(encoder))
	require.EqualValues(t, "23e629f9aa6b2c46aa8fa836770e7a7a5f0561646d696e027b7d9b2a645b2100000000010ce803", hex.EncodeToString(b.Bytes()))
}

func TestCreateGameOperation_SerializationManyMarkets(t *testing.T) {
	time := time.Unix(1461605400, 0)

	uuid := uuid.UUID{}
	require.NoError(t, uuid.UnmarshalText([]byte("e629f9aa-6b2c-46aa-8fa8-36770e7a7a5f")))

	op := CreateGameOperation{
		UUID:                uuid,
		Moderator:           "moderator_name",
		JsonMetadata:        "{}",
		StartTime:           Time{&time},
		AutoResolveDelaySec: 33,
		Markets: []Market{
			Market{&YesNoMarket{
				ID: MarketResultHome,
			}},
			Market{&YesNoMarket{
				ID: MarketResultDraw,
			}},
			Market{&YesNoMarket{
				ID: MarketResultAway,
			}},
			Market{&YesNoMarket{
				ID: MarketRoundHome,
			}},
			Market{&OverUnderMarket{
				ID:        MarketHandicap,
				Threshold: -500,
			}},
			Market{&OverUnderMarket{
				ID:        MarketHandicap,
				Threshold: 0,
			}},
			Market{&OverUnderMarket{
				ID:        MarketHandicap,
				Threshold: 1000,
			}},
			Market{&YesNoMarket{
				ID: MarketCorrectScoreHome,
			}},
			Market{&YesNoMarket{
				ID: MarketCorrectScoreDraw,
			}},
			Market{&YesNoMarket{
				ID: MarketCorrectScoreAway,
			}},
			Market{&ScoreYesNoMarket{
				ID:   MarketCorrectScore,
				Home: 1,
				Away: 0,
			}},
			Market{&ScoreYesNoMarket{
				ID:   MarketCorrectScore,
				Home: 1,
				Away: 1,
			}},
			Market{&YesNoMarket{
				ID: MarketGoalHome,
			}},
			Market{&YesNoMarket{
				ID: MarketGoalBoth,
			}},
			Market{&YesNoMarket{
				ID: MarketGoalAway,
			}},
			Market{&OverUnderMarket{
				ID:        MarketTotal,
				Threshold: 0,
			}},
			Market{&OverUnderMarket{
				ID:        MarketTotal,
				Threshold: 500,
			}},
			Market{&OverUnderMarket{
				ID:        MarketTotal,
				Threshold: 1000,
			}},
		},
	}

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)
	require.NoError(t, op.MarshalTransaction(encoder))
	require.EqualValues(
		t,
		"23e629f9aa6b2c46aa8fa836770e7a7a5f0e6d6f64657261746f725f6e616d65027b7d18541e5721000000001200010203040cfe04000004e80305060708010000000801000100090a0b0c00000cf4010ce803",
		hex.EncodeToString(b.Bytes()))
}

func TestPostGameResultsOperation_Serialization(t *testing.T) {
	uuid := uuid.UUID{}
	require.NoError(t, uuid.UnmarshalText([]byte("e629f9aa-6b2c-46aa-8fa8-36770e7a7a5f")))

	op := PostGameResultsOperation{
		UUID:      uuid,
		Moderator: "homer",
		Wincases: []Wincase{
			Wincase{
				&YesNoWincase{
					ID: WincaseResultHomeYes,
				}},
			Wincase{
				&YesNoWincase{
					ID: WincaseResultDrawNo,
				}},
			Wincase{
				&YesNoWincase{
					ID: WincaseResultAwayYes,
				}},
			Wincase{
				&YesNoWincase{
					ID: WincaseRoundHomeNo,
				}},
			Wincase{
				&OverUnderWincase{
					ID:        WincaseHandicapOver,
					Threshold: 1000,
				}},
			Wincase{
				&OverUnderWincase{
					ID:        WincaseHandicapUnder,
					Threshold: -500,
				}},
			Wincase{
				&OverUnderWincase{
					ID:        WincaseHandicapUnder,
					Threshold: 0,
				}},
			Wincase{
				&YesNoWincase{
					ID: WincaseCorrectScoreHomeYes,
				}},
			Wincase{
				&YesNoWincase{
					ID: WincaseCorrectScoreDrawNo,
				}},
			Wincase{
				&YesNoWincase{
					ID: WincaseCorrectScoreAwayNo,
				}},
			Wincase{
				&ScoreYesNoWincase{
					ID:   WincaseCorrectScoreYes,
					Home: 1,
					Away: 2,
				}},
			Wincase{
				&ScoreYesNoWincase{
					ID:   WincaseCorrectScoreNo,
					Home: 3,
					Away: 2,
				}},
			Wincase{
				&YesNoWincase{
					ID: WincaseGoalHomeYes,
				}},
			Wincase{
				&YesNoWincase{
					ID: WincaseGoalBothNo,
				}},
			Wincase{
				&YesNoWincase{
					ID: WincaseGoalAwayYes,
				}},
			Wincase{
				&OverUnderWincase{
					ID:        WincaseTotalOver,
					Threshold: 0,
				}},
			Wincase{
				&OverUnderWincase{
					ID:        WincaseTotalUnder,
					Threshold: 1000,
				}},
		},
	}
	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)

	require.NoError(t, op.MarshalTransaction(encoder))

	require.EqualValues(
		t,
		"27e629f9aa6b2c46aa8fa836770e7a7a5f05686f6d6572110003040708e803090cfe0900000a0d0f1001000200110300020012151618000019e803",
		hex.EncodeToString(b.Bytes()))
}

func TestPostBetOperation_Serialization(t *testing.T) {
	uuid := uuid.UUID{}
	require.NoError(t, uuid.UnmarshalText([]byte("e629f9aa-6b2c-46aa-8fa8-36770e7a7a5f")))

	asset, err := AssetFromString("10.000000000 SCR")
	require.NoError(t, err)

	op := PostBetOperation{
		UUID:     uuid,
		Better:   "admin",
		GameUUID: uuid,
		Wincase: Wincase{
			&ScoreYesNoWincase{
				ID:   WincaseCorrectScoreYes,
				Home: 17,
				Away: 23,
			}},
		Odds: Odds{
			Numerator:   1,
			Denominator: 2,
		},
		Stake: *asset,
		Live:  true,
	}
	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)

	require.NoError(t, op.MarshalTransaction(encoder))

	require.EqualValues(
		t,
		"28e629f9aa6b2c46aa8fa836770e7a7a5f0561646d696ee629f9aa6b2c46aa8fa836770e7a7a5f1011001700010000000200000000e40b5402000000095343520000000001",
		hex.EncodeToString(b.Bytes()))
}

func TestCancelPendingBetsOperation_Serialization(t *testing.T) {
	uuids := make([]uuid.UUID, 0)

	var id uuid.UUID
	require.NoError(t, id.UnmarshalText([]byte("e629f9aa-6b2c-46aa-8fa8-36770e7a7a5f")))

	uuids = append(uuids, id)

	op := CancelPendingBetsOperation{
		BetIDs: uuids,
		Better: "admin",
	}

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)

	require.NoError(t, op.MarshalTransaction(encoder))

	require.EqualValues(
		t,
		"2901e629f9aa6b2c46aa8fa836770e7a7a5f0561646d696e",
		hex.EncodeToString(b.Bytes()))
}

func TestCreateGameOperation_UnmarshalJSON(t *testing.T) {
	testJson := `{
                                              "uuid":"e629f9aa-6b2c-46aa-8fa8-36770e7a7a5f",
                                              "moderator":"daddy",
                                              "json_metadata":"{}",
                                              "start_time":"1970-01-01T00:00:00",
                                              "auto_resolve_delay_sec":33,
                                              "game":[
                                                 "soccer_game",
                                                 {}
                                              ],
                                              "markets":[
                                                 [
                                                    "correct_score_home",
                                                    {}
                                                 ],
                                                 [
                                                    "correct_score",
                                                    {
                                                       "home":1,
                                                       "away":2
                                                    }
                                                 ]
                                              ]
                                           }`
	var game CreateGameOperation
	require.NoError(t, json.Unmarshal([]byte(testJson), &game))
	require.EqualValues(t, "daddy", game.Moderator)
	require.EqualValues(t, "{}", game.JsonMetadata)
	require.EqualValues(t, 33, game.AutoResolveDelaySec)
	require.Len(t, game.Markets, 2)
	require.IsType(t, &YesNoMarket{}, game.Markets[0].MarketInterface)
	require.EqualValues(t, 5, game.Markets[0].MarketInterface.(*YesNoMarket).ID)
	require.IsType(t, &ScoreYesNoMarket{}, game.Markets[1].MarketInterface)
	require.EqualValues(t, game.Markets[1].MarketInterface.(*ScoreYesNoMarket).Home, 1)
	require.EqualValues(t, game.Markets[1].MarketInterface.(*ScoreYesNoMarket).Away, 2)
	require.EqualValues(t, 8, game.Markets[1].MarketInterface.(*ScoreYesNoMarket).ID)
}

func Test_CreateNFTOperation_Marshall(t *testing.T) {
	op := CreateNFTOperation{
		OwnerAccount: "operator",
		UUID:         uuid.MustParse("aa3b2bdc-176e-5e8a-9b48-7dba3aa10044"),
		Name:         "rocket",
		JSONMetadata: "{}",
		InitialPower: 100,
	}

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)

	require.NoError(t, op.MarshalTransaction(encoder))

	require.EqualValues(t,
		"2b086f70657261746f72aa3b2bdc176e5e8a9b487dba3aa1004406726f636b6574027b7d64000000",
		hex.EncodeToString(b.Bytes()),
	)
}

func Test_UpdateNFTMetadataOperation_Marshall(t *testing.T) {
	op := UpdateNFTMetadataOperation{
		Moderator:    "operator",
		UUID:         uuid.MustParse("aa3b2bdc-176e-5e8a-9b48-7dba3aa10044"),
		JSONMetadata: "{}",
	}

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)

	require.NoError(t, op.MarshalTransaction(encoder))

	require.EqualValues(t,
		"2c086f70657261746f72aa3b2bdc176e5e8a9b487dba3aa10044027b7d",
		hex.EncodeToString(b.Bytes()),
	)
}

func Test_CreateGameRoundOperation_Marshall(t *testing.T) {
	op := CreateGameRoundOperation{
		Owner:           "operator",
		UUID:            uuid.MustParse("aa3b2bdc-176e-5e8a-9b48-7dba3aa10044"),
		VerificationKey: "038b2fbf4e4f066f309991b9c30cb8f887853e54c76dc705f5ece736ead6c856",
		Seed:            "052d8bec6d55f8c489e837c19fa372e00e3433ebfe0068af658e36b0bd1eb722",
	}

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)

	require.NoError(t, op.MarshalTransaction(encoder))

	require.Equal(t,
		"2d086f70657261746f72aa3b2bdc176e5e8a9b487dba3aa1004440303338623266626634653466303636663330393939316239633330636238663838373835336535346337366463373035663565636537333665616436633835364030353264386265633664353566386334383965383337633139666133373265303065333433336562666530303638616636353865333662306264316562373232",
		hex.EncodeToString(b.Bytes()),
	)
}

func Test_UpdateGameRoundResultOperation_Marshall(t *testing.T) {
	op := UpdateGameRoundResultOperation{
		Owner:  "operator",
		UUID:   uuid.MustParse("aa3b2bdc-176e-5e8a-9b48-7dba3aa10044"),
		Proof:  "638f675cd4313ae84aede4940b7691acd904dec141e444187dcec59f2a25a7a4ef5aa2fe3f88cf235c0d63aa6935bef69d5b70caca0d9b4028f75121d030f80a5a4bcf97b36a868ea9a4c2aaa9013200",
		Vrf:    "6a196a14e4f9fce66112b1b7ac98f2bcd73352b918d298e0b9f894519a65202dba03ddaa5183190bf2b5cd551f9ef14c8d8b02cf15d0188bbc9bcc6a80d7f91c",
		Result: 100,
	}

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)

	require.NoError(t, op.MarshalTransaction(encoder))

	require.Equal(t,
		"2e086f70657261746f72aa3b2bdc176e5e8a9b487dba3aa10044a001363338663637356364343331336165383461656465343934306237363931616364393034646563313431653434343138376463656335396632613235613761346566356161326665336638386366323335633064363361613639333562656636396435623730636163613064396234303238663735313231643033306638306135613462636639376233366138363865613961346332616161393031333230308001366131393661313465346639666365363631313262316237616339386632626364373333353262393138643239386530623966383934353139613635323032646261303364646161353138333139306266326235636435353166396566313463386438623032636631356430313838626263396263633661383064376639316364000000",
		hex.EncodeToString(b.Bytes()),
	)
}

func TestAccountCreateByCommitteeOperation_MarshalTransaction(t *testing.T) {
	op := AccountCreateByCommitteeOperation{
		Creator:        "alice",
		NewAccountName: "bob",
		Owner: Authority{
			WeightThreshold: 1,
			AccountAuths:    NewAccountAuthorityMap(),
			KeyAuths:        NewKeyAuthorityMap(KeyAuthority{Key: "SCR7zPNg5nAsJjP9gvMfQ4UnAwDwf91WPYC8KFzobtMuQ52ns1D6T", Weight: 1}),
		},
		Active: Authority{
			WeightThreshold: 1,
			AccountAuths:    NewAccountAuthorityMap(),
			KeyAuths:        NewKeyAuthorityMap(KeyAuthority{Key: "SCR7SHdKpjpWyfj32tGQBeijFfokmCARjKSBynqDwN1ZAbQRW5rWa", Weight: 1}),
		},
		Posting: Authority{
			WeightThreshold: 1,
			AccountAuths:    NewAccountAuthorityMap(),
			KeyAuths:        NewKeyAuthorityMap(KeyAuthority{Key: "SCR5jPZF7PMgTpLqkdfpMu8kXea8Gio6E646aYpTgcjr9qMLrAgnL", Weight: 1}),
		},
		MemoKey:      "SCR5jPZF7PMgTpLqkdfpMu8kXea8Gio6E646aYpTgcjr9qMLrAgnL",
		JsonMetadata: "",
	}

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)

	require.NoError(t, op.MarshalTransaction(encoder))

	require.Equal(t,
		"0505616c69636503626f6201000000000103987a5a967458c114c15091198c06a822f54b494ea486204551a53f85effa31420100010000000001034f97d09e6de4778300ed176403e5b4298bfd62f0fb6edb4a6072e7214318d9030100010000000001026f0896f24d94252c351715bfe6052bbf9ea820e805bd47c2496c626d3467da5d0100026f0896f24d94252c351715bfe6052bbf9ea820e805bd47c2496c626d3467da5d00",
		hex.EncodeToString(b.Bytes()),
	)
}

func TestAccountCreateOperation(t *testing.T) {
	op := AccountCreateOperation{
		Fee:            *AssetFromFloat(0.000000750),
		Creator:        "alice",
		NewAccountName: "bob",
		Owner: Authority{
			WeightThreshold: 1,
			AccountAuths:    NewAccountAuthorityMap(),
			KeyAuths:        NewKeyAuthorityMap(KeyAuthority{Key: "SCR7zPNg5nAsJjP9gvMfQ4UnAwDwf91WPYC8KFzobtMuQ52ns1D6T", Weight: 1}),
		},
		Active: Authority{
			WeightThreshold: 1,
			AccountAuths:    NewAccountAuthorityMap(),
			KeyAuths:        NewKeyAuthorityMap(KeyAuthority{Key: "SCR7SHdKpjpWyfj32tGQBeijFfokmCARjKSBynqDwN1ZAbQRW5rWa", Weight: 1}),
		},
		Posting: Authority{
			WeightThreshold: 1,
			AccountAuths:    NewAccountAuthorityMap(),
			KeyAuths:        NewKeyAuthorityMap(KeyAuthority{Key: "SCR5jPZF7PMgTpLqkdfpMu8kXea8Gio6E646aYpTgcjr9qMLrAgnL", Weight: 1}),
		},
		MemoKey:      "SCR5jPZF7PMgTpLqkdfpMu8kXea8Gio6E646aYpTgcjr9qMLrAgnL",
		JsonMetadata: "",
	}

	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)

	require.NoError(t, op.MarshalTransaction(encoder))

	require.Equal(t,
		"06ee02000000000000095343520000000005616c69636503626f6201000000000103987a5a967458c114c15091198c06a822f54b494ea486204551a53f85effa31420100010000000001034f97d09e6de4778300ed176403e5b4298bfd62f0fb6edb4a6072e7214318d9030100010000000001026f0896f24d94252c351715bfe6052bbf9ea820e805bd47c2496c626d3467da5d0100026f0896f24d94252c351715bfe6052bbf9ea820e805bd47c2496c626d3467da5d00",
		hex.EncodeToString(b.Bytes()),
	)
}
