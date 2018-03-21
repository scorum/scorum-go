package types

type Authority struct {
	WeightThreshold int32          `json:"weight_threshold"`
	AccountAuths    StringInt64Map `json:"account_auths"`
	KeyAuths        StringInt64Map `json:"key_auths"`
}
