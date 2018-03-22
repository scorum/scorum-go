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
	CommentOptions,
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
	AuthorReward,
	CommentBenefactorReward,
	CommentPayoutUpdate,
	CommentReward,
	CurationReward,
	FillScorumpowerWithdraw,
	Hardfork,
	ProducerRewardOperation,
	ReturnScorumpowerDelegation,
	ShutdownWitness,
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
	CommentOptions                       OpType = "comment_options"
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
	AuthorReward                         OpType = "author_reward"
	CommentBenefactorReward              OpType = "comment_benefactor_reward"
	CommentPayoutUpdate                  OpType = "comment_payout_update"
	CommentReward                        OpType = "comment_reward"
	CurationReward                       OpType = "curation_reward"
	FillScorumpowerWithdraw              OpType = "fill_scorumpower_withdraw"
	Hardfork                             OpType = "hardfork"
	ProducerRewardOperation              OpType = "producer_reward_operation"
	ReturnScorumpowerDelegation          OpType = "return_scorumpower_delegation"
	ShutdownWitness                      OpType = "shutdown_witness"
)
