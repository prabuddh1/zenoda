package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
	HasAccount(context.Context, sdk.AccAddress) bool
	SetAccount(context.Context, sdk.AccountI)
	NewAccountWithAddress(context.Context, sdk.AccAddress) sdk.AccountI
	// GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetSupply(ctx context.Context, denom string) sdk.Coin
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
