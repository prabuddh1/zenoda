package rewards

import (
	"fmt"

	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"zenoda/x/rewards/keeper"
	"zenoda/x/rewards/types"
)

// InitGenesis initializes the module's state from the provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	ctx.Logger().Info("🚀 Initializing rewards module genesis...")

	// // Get the rewards module account address
	// moduleAddr := k.GetAccountKeeper().GetModuleAddress(types.ModuleName)

	// // Retrieve the account from the store using the GetAccount method
	// account := k.GetAccountKeeper().GetAccount(ctx, moduleAddr)

	// // Check if the account exists
	// if account == nil {
	// 	ctx.Logger().Info("❌ Rewards module account does not exist, creating...")
	// 	// Create and set the rewards module account
	// 	account = k.GetAccountKeeper().NewAccountWithAddress(ctx, moduleAddr)
	// 	k.GetAccountKeeper().SetAccount(ctx, account) // Set the newly created account
	// 	ctx.Logger().Info("✅ Created rewards module account", "address", moduleAddr.String())
	// } else {
	// 	ctx.Logger().Info("✅ Rewards module account already exists", "address", moduleAddr.String())
	// }

	// // Check again to ensure the module account is registered before minting
	// if !k.GetAccountKeeper().HasAccount(ctx, moduleAddr) {
	// 	panic("module account rewards does not exist after creation")
	// }

	// // Verify the module account address
	// ctx.Logger().Info("Module account address for minting", "address", moduleAddr.String())

	// Directly create a new module account using the module name "rewards" (from keys.go)
	moduleAddr := sdk.AccAddress(types.ModuleName) // "rewards" should be used from keys.go
	// Create the new module account using NewAccountWithAddress
	account := k.GetAccountKeeper().NewAccountWithAddress(ctx, moduleAddr)

	// Set the newly created account in the store
	k.GetAccountKeeper().SetAccount(ctx, account)

	// Log the creation of the module account
	ctx.Logger().Info("✅ Created rewards module account", moduleAddr)
	ctx.Logger().Info("✅ Created rewards module account", "address", moduleAddr.String())

	// Retrieve the account from the store using GetAccount
	account = k.GetAccountKeeper().GetAccount(ctx, moduleAddr)

	// Log account details to verify
	ctx.Logger().Info("Account retrieved", "address", moduleAddr.String(), "account details", account)

	// Define the initial amount each predefined wallet will receive
	initialAmount := sdk.NewCoin(types.EGVDenom, math.NewInt(1000)) // Each gets 1000 EGV
	totalSupply := math.NewInt(int64(len(genState.Params.PredefinedWallets))).Mul(initialAmount.Amount)

	// Ensure that the module account exists before minting
	ctx.Logger().Info("Checking if the module account exists before minting", "moduleAccountExists", k.GetAccountKeeper().HasAccount(ctx, moduleAddr))

	// Explicitly log the module address before minting
	ctx.Logger().Info("Attempting to mint coins for the module account", "moduleAddress", moduleAddr.String())

	// Mint tokens only once for the total supply
	err := k.GetBankKeeper().MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(types.EGVDenom, totalSupply)))
	if err != nil {
		ctx.Logger().Error("❌ Failed to mint initial EGV tokens", "error", err)
		panic(err)
	}

	ctx.Logger().Info("💰 Minted total EGV tokens for initial distribution", "amount", totalSupply)

	// Distribute to predefined wallets
	for _, addressStr := range genState.Params.PredefinedWallets {
		address, err := sdk.AccAddressFromBech32(addressStr)
		if err != nil {
			ctx.Logger().Error("❌ Invalid predefined wallet address", "address", addressStr, "error", err)
			panic(fmt.Sprintf("invalid address in genesis: %s", err))
		}

		// Send tokens from module account to predefined wallet
		err = k.GetBankKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, sdk.NewCoins(initialAmount))
		if err != nil {
			ctx.Logger().Error("❌ Failed to send EGV tokens to address", "address", addressStr, "error", err)
			panic(err)
		}

		ctx.Logger().Info("✅ Distributed initial EGV tokens", "address", addressStr, "amount", initialAmount.Amount)
	}

	// Store total supply in the keeper
	k.SetTotalSupply(ctx, sdk.NewCoin(types.EGVDenom, totalSupply))

	ctx.Logger().Info("✅ Rewards module genesis successfully initialized")
}

// ExportGenesis exports the module's state.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	return genesis
}
