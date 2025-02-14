package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	math "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"zenoda/x/rewards/types"
	//"github.com/cosmos/cosmos-sdk/x/auth/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		authority     string
		bankKeeper    types.BankKeeper
		accountKeeper types.AccountKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:           cdc,
		storeService:  storeService,
		authority:     authority,
		logger:        logger,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
	}
}

// ---------------------- CORE UTILS ----------------------

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetBankKeeper returns the bank keeper instance
func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

// Getter for AccountKeeper
func (k Keeper) GetAccountKeeper() types.AccountKeeper {
	return k.accountKeeper
}

// SetTotalSupply stores the total supply of EGV tokens
func (k Keeper) SetTotalSupply(ctx sdk.Context, supply sdk.Coin) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&supply)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal total supply: %s", err))
	}
	_ = store.Set([]byte(types.TotalSupplyKey), bz)
}

// ---------------------- PARAMETER ACCESS ----------------------

// GetParams fetches the module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := k.storeService.OpenKVStore(ctx)
	var params types.Params
	bz, err := store.Get([]byte(types.ParamsKey))
	if err != nil || bz == nil {
		panic("failed to retrieve params from store")
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// ---------------------- ADDRESS MANAGEMENT ----------------------

// GetPredefinedAddresses fetches the predefined governance addresses from the module parameters.
func (k Keeper) GetPredefinedAddresses(ctx sdk.Context) []sdk.AccAddress {
	params := k.GetParams(ctx)
	var addresses []sdk.AccAddress

	for _, addrStr := range params.PredefinedWallets {
		accAddr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			ctx.Logger().Error("Invalid predefined address", "address", addrStr, "error", err)
			continue
		}
		addresses = append(addresses, accAddr)
	}
	return addresses
}

// ---------------------- TRANSACTION COUNT TRACKING ----------------------

// Increment transaction count for a given address
func (k Keeper) IncrementTransactionCount(ctx sdk.Context, addr sdk.AccAddress) {
	store := k.storeService.OpenKVStore(ctx)

	// Check if the address is part of the predefined set
	predefinedAddresses := k.GetPredefinedAddresses(ctx)
	if !k.isPredefinedAddress(addr, predefinedAddresses) {
		return
	}

	addrKey := append([]byte(types.TransactionCountKey), addr.Bytes()...)

	// Get current count
	bz, err := store.Get(addrKey)
	var count uint64
	if err == nil && bz != nil {
		count = sdk.BigEndianToUint64(bz)
	}

	// Increment count
	count++

	// Store updated count
	_ = store.Set(addrKey, sdk.Uint64ToBigEndian(count))

	// Increment total network transactions (includes all network addresses)
	k.IncrementTotalTransactions(ctx)
}

// Helper function to check if an address is predefined
func (k Keeper) isPredefinedAddress(addr sdk.AccAddress, predefined []sdk.AccAddress) bool {
	for _, predefinedAddr := range predefined {
		if predefinedAddr.Equals(addr) {
			return true
		}
	}
	return false
}

// Get transaction count for an address
func (k Keeper) GetTransactionCount(ctx sdk.Context, addr sdk.AccAddress) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	addrKey := append([]byte(types.TransactionCountKey), addr.Bytes()...)

	bz, err := store.Get(addrKey)
	if err != nil || bz == nil {
		return 0 // Default to 0 if not found
	}
	return sdk.BigEndianToUint64(bz)
}

// ---------------------- TOTAL TRANSACTION TRACKING ----------------------

// Increment total transaction count in the network
func (k Keeper) IncrementTotalTransactions(ctx sdk.Context) {
	store := k.storeService.OpenKVStore(ctx)
	totalTxKey := []byte(types.TotalTxKey)

	// Get current count
	bz, err := store.Get(totalTxKey)
	var total uint64
	if err == nil && bz != nil {
		total = sdk.BigEndianToUint64(bz)
	}

	// Increment count
	total++

	// Store updated total transactions
	_ = store.Set(totalTxKey, sdk.Uint64ToBigEndian(total))
}

// Get total transactions in the network
func (k Keeper) GetTotalTransactions(ctx sdk.Context) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	totalTxKey := []byte(types.TotalTxKey)

	bz, err := store.Get(totalTxKey)
	if err != nil || bz == nil {
		return 0 // Default to 0 if not found
	}
	return sdk.BigEndianToUint64(bz)
}

// ---------------------- EGV SUPPLY AND INFLATION ----------------------

// GetTotalSupply returns the total supply of EGV tokens
func (k Keeper) GetTotalSupply(ctx sdk.Context) sdk.Coin {
	return k.bankKeeper.GetSupply(ctx, types.EGVDenom)
}

// GetInflationRate returns the current inflation rate for EGV
func (k Keeper) GetInflationRate(ctx sdk.Context) math.LegacyDec {
	params := k.GetParams(ctx)
	inflationRate, err := params.GetInflationRateAsDec()
	if err != nil {
		panic(fmt.Sprintf("invalid inflation rate stored in params: %s", err))
	}
	return inflationRate
}

// Lock setter: Sets a lock in the KV store.
func (k Keeper) SetLock(ctx sdk.Context, lockKey string) {
	store := k.storeService.OpenKVStore(ctx)
	err := store.Set([]byte(lockKey), []byte("locked"))
	if err != nil {
		// Handle error (log it or return if you want to fail early)
		k.Logger().Error("Error setting lock", "error", err)
		panic("Failed to set lock")
	}
}

// Lock checker: Checks if the lock is set.
func (k Keeper) IsLockSet(ctx sdk.Context, lockKey string) bool {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get([]byte(lockKey))
	if err != nil {
		// Handle error (log it or return false if you want to allow continuing)
		k.Logger().Error("Error retrieving lock", "error", err)
		return false
	}
	return bz != nil && string(bz) == "locked"
}

// Clear lock: Removes the lock after operation is complete.
func (k Keeper) ClearLock(ctx sdk.Context, lockKey string) {
	store := k.storeService.OpenKVStore(ctx)
	err := store.Delete([]byte(lockKey))
	if err != nil {
		// Handle error (log it or return if you want to fail early)
		k.Logger().Error("Error clearing lock", "error", err)
		panic("Failed to clear lock")
	}
}

// func (k Keeper) MintEGV(ctx sdk.Context, amount sdk.Coin) error {
// 	// Mint coins to the module account
// 	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(amount)); err != nil {
// 		return err
// 	}

// 	// Send the minted coins to the rewards pool module account
// 	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.RewardsModuleName, sdk.NewCoins(amount))
// }

// // SendInitialEGV sends the initial EGV amount to a given address
// func (k Keeper) SendInitialEGV(ctx sdk.Context, addr sdk.AccAddress, coin sdk.Coin) error {
// 	// Use the BankKeeper to get the supply (total amount) of EGV in the module
// 	supply := k.bankKeeper.GetSupply(ctx, types.EGVDenom)

// 	// Check if the module account has enough balance (checking the supply in this case)
// 	if supply.IsLT(coin) {
// 		return fmt.Errorf("insufficient balance in the rewards module to send initial EGV to %s: available %s, needed %s", addr, supply.String(), coin.String())
// 	}

// 	// Use BankKeeper to send coins from the rewards module to the predefined wallet
// 	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(coin)); err != nil {
// 		return fmt.Errorf("failed to send initial EGV to address %s: %v", addr, err)
// 	}

// 	return nil
// }
