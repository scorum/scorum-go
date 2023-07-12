package types

type OpType string

// Code returns the operation code associated with the given operation type.
func (kind OpType) Code() uint16 {
	return opCodes[kind]
}

// opCodes keeps mapping operation type -> operation code.
var opCodes map[OpType]uint16

func init() {
	opCodes = make(map[OpType]uint16, len(opTypes))
	for i, opType := range opTypes {
		opCodes[opType] = uint16(i)
	}
}

var opTypes = []OpType{
	VoteOpType,
	CommentOpType,
	TransferOpType,
	TransferToScorumpowerOpType,
	WithdrawScorumpowerOpType,
	AccountCreateByCommitteeOpType,
	AccountCreateOpType,
	AccountCreateWithDelegationOpType,
	AccountUpdateOpType,
	WitnessUpdateOpType,
	AccountWitnessVoteOpType,
	AccountWitnessProxyOpType,
	DeleteCommentOpType,
	CommentOptionsOpType,
	SetWithdrawScorumpowerRouteToAccount,
	SetWithdrawScorumpowerRouteToDevPool,
	ProveAuthority,
	RequestAccountRecovery,
	RecoverAccount,
	ChangeRecoveryAccount,
	EscrowApprove,
	EscrowDispute,
	EscrowRelease,
	EscrowTransfer,
	DeclineVotingRights,
	DelegateScorumpower,
	CreateBudget,
	CloseBudget,
	ProposalVoteOperation,
	ProposalCreateOperation,
	AtomicswapInitiateOperation,
	AtomicswapRedeemOperation,
	AtomicswapRefundOperation,

	CloseBudgetByAdvertisingModeratorOperation,
	UpdateBudgetOperation,

	CreateGame,
	CancelGame,
	UpdateGameMarkets,
	UpdateGameStartTime,
	PostGameResults,
	PostBet,
	CancelPendingBets,

	DelegateSPFromRegPool,

	CreateNFT,
	UpdateNFTMetadata,
	CreateGameRound,
	UpdateGameRoundResult,
	AdjustNFTExperience,
	UpdateNFTName,

	BurnOperationOpType,

	// virtual operations
	CommentBenefactorReward,
	CommentPayoutUpdate,
	CommentReward,
	CurationReward,
	FillScorumpowerWithdraw,
	Hardfork,
	ProducerRewardOpType,
	ReturnScorumpowerDelegation,
	ShutdownWitness,
	WitnessMissBlock,
	ExpiredContractRefund,
	AccFinishedVestingWithdraw,
	DevpoolFinishedVestingWithdraw,
	AccToAccVestingWithdraw,
	DevpoolToAccVestingWithdraw,
	AccToDevpoolVestingWithdraw,
	DevpoolToDevpoolVesting,
	ProposalVirtual,
	ActiveSpHoldersRewardLegacy,
	AllocateCashFromAdvertisingBudget,
	CashBackFromAdvertisingBudgetToOwner,
	ClosingBudget,
	BetsMatched,
	GameStatusChanged,
	BetResolved,
	BetCancelled,
	BetRestored,
	BetUpdated,
}

const (
	VoteOpType                           OpType = "vote"
	CommentOpType                        OpType = "comment"
	TransferOpType                       OpType = "transfer"
	TransferToScorumpowerOpType          OpType = "transfer_to_scorumpower"
	WithdrawScorumpowerOpType            OpType = "withdraw_scorumpower"
	AccountCreateByCommitteeOpType       OpType = "account_create_by_committee"
	AccountCreateOpType                  OpType = "account_create"
	AccountCreateWithDelegationOpType    OpType = "account_create_with_delegation"
	AccountUpdateOpType                  OpType = "account_update"
	WitnessUpdateOpType                  OpType = "witness_update"
	AccountWitnessVoteOpType             OpType = "account_witness_vote"
	AccountWitnessProxyOpType            OpType = "account_witness_proxy"
	DeleteCommentOpType                  OpType = "delete_comment"
	CommentOptionsOpType                 OpType = "comment_options"
	SetWithdrawScorumpowerRouteToAccount OpType = "set_withdraw_scorumpower_route_to_account"
	SetWithdrawScorumpowerRouteToDevPool OpType = "set_withdraw_scorumpower_route_to_dev_pool"
	ProveAuthority                       OpType = "prove_authority"
	RequestAccountRecovery               OpType = "request_account_recovery"
	RecoverAccount                       OpType = "recover_account"
	ChangeRecoveryAccount                OpType = "change_recovery_account"
	EscrowApprove                        OpType = "escrow_approve"
	EscrowDispute                        OpType = "escrow_dispute"
	EscrowRelease                        OpType = "escrow_release"
	EscrowTransfer                       OpType = "escrow_transfer"
	DeclineVotingRights                  OpType = "decline_voting_rights"
	DelegateScorumpower                  OpType = "delegate_scorumpower"
	CreateBudget                         OpType = "create_budget"
	CloseBudget                          OpType = "close_budget"
	ProposalVoteOperation                OpType = "proposal_vote_operation"
	ProposalCreateOperation              OpType = "proposal_create_operation"
	AtomicswapInitiateOperation          OpType = "atomicswap_initiate_operation"
	AtomicswapRedeemOperation            OpType = "atomicswap_redeem_operation"
	AtomicswapRefundOperation            OpType = "atomicswap_refund_operation"
	BurnOperationOpType                  OpType = "burn"

	CloseBudgetByAdvertisingModeratorOperation OpType = "close_budget_by_advertising_moderator"
	UpdateBudgetOperation                      OpType = "update_budget"

	CreateGame          OpType = "create_game"
	CancelGame          OpType = "cancel_game"
	UpdateGameMarkets   OpType = "update_game_markets"
	UpdateGameStartTime OpType = "update_game_start_time"
	PostGameResults     OpType = "post_game_results"
	PostBet             OpType = "post_bet"
	CancelPendingBets   OpType = "cancel_pending_bets"

	DelegateSPFromRegPool OpType = "delegate_sp_from_reg_pool"

	CreateNFT             OpType = "create_nft"
	UpdateNFTMetadata     OpType = "update_nft_meta"
	CreateGameRound       OpType = "create_game_round"
	UpdateGameRoundResult OpType = "update_game_round_result"
	AdjustNFTExperience   OpType = "adjust_nft_experience"
	UpdateNFTName         OpType = "update_nft_name"

	// virtual operations
	AuthorReward                OpType = "author_reward"
	CommentBenefactorReward     OpType = "comment_benefactor_reward"
	CommentPayoutUpdate         OpType = "comment_payout_update"
	CommentReward               OpType = "comment_reward"
	CurationReward              OpType = "curation_reward"
	FillScorumpowerWithdraw     OpType = "fill_scorumpower_withdraw"
	Hardfork                    OpType = "hardfork"
	ProducerRewardOpType        OpType = "producer_reward"
	ReturnScorumpowerDelegation OpType = "return_scorumpower_delegation"
	ShutdownWitness             OpType = "shutdown_witness"

	WitnessMissBlock                     OpType = "witness_miss_block"
	ExpiredContractRefund                OpType = "expired_contract_refund"
	AccFinishedVestingWithdraw           OpType = "acc_finished_vesting_withdraw"
	DevpoolFinishedVestingWithdraw       OpType = "devpool_finished_vesting_withdraw"
	AccToAccVestingWithdraw              OpType = "acc_to_acc_vesting_withdraw"
	DevpoolToAccVestingWithdraw          OpType = "devpool_to_acc_vesting_withdraw"
	AccToDevpoolVestingWithdraw          OpType = "acc_to_devpool_vesting_withdraw"
	DevpoolToDevpoolVesting              OpType = "devpool_to_devpool_vesting_withdraw"
	ProposalVirtual                      OpType = "proposal_virtual"
	ActiveSpHoldersRewardLegacy          OpType = "active_sp_holders_reward_legacy"
	AllocateCashFromAdvertisingBudget    OpType = "allocate_cash_from_advertising_budget"
	CashBackFromAdvertisingBudgetToOwner OpType = "cash_back_from_advertising_budget_to_owner"
	ClosingBudget                        OpType = "closing_budget"
	BetsMatched                          OpType = "bets_matched"
	GameStatusChanged                    OpType = "game_status_changed"

	BetResolved  OpType = "bet_resolved"
	BetCancelled OpType = "bet_cancelled"
	BetRestored  OpType = "bet_restored"
	BetUpdated   OpType = "bet_updated"
)
