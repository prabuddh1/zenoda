package types

const (
	// ModuleName defines the module name
	ModuleName = "rewards"

	// RewardsModuleName defines the name of the rewards pool module account
	RewardsModuleName = "rewards_pool"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_rewards"

	// TransactionCountKey is the prefix to store transaction count per address
	TransactionCountKey = "transaction_count"

	// TotalTxKey is the key for storing total transactions in the network
	TotalTxKey = "total_transactions"

	// EGVDenom is the denomination for the EGV token
	EGVDenom = "egv"

	TotalSupplyKey = "TotalSupply"
)

var (
	ParamsKey = []byte("p_rewards")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
