package rewards

import (
	"fmt"

	"zenoda/x/rewards/keeper"
	"zenoda/x/rewards/types"

	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set the module parameters from the genesis state
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	// Create the rewards module account (rewards)
	rewardsModuleAccountAddr := k.GetAccountKeeper().GetModuleAddress(types.ModuleName)            // "rewards"
	fmt.Println("Checking if rewards module account exists at address:", rewardsModuleAccountAddr) // Log this address

	if !k.GetAccountKeeper().HasAccount(ctx, rewardsModuleAccountAddr) {
		fmt.Println("Creating Rewards Module Account (rewards) at address:", rewardsModuleAccountAddr)
		account := k.GetAccountKeeper().NewAccountWithAddress(ctx, rewardsModuleAccountAddr)
		k.GetAccountKeeper().SetAccount(ctx, account)
	} else {
		fmt.Println("Rewards Module Account (rewards) already exists at address:", rewardsModuleAccountAddr)
	}

	// Create the rewards pool module account (rewards_pool)
	rewardsPoolAccountAddr := k.GetAccountKeeper().GetModuleAddress(types.RewardsModuleName)          // "rewards_pool"
	fmt.Println("Checking if rewards pool module account exists at address:", rewardsPoolAccountAddr) // Log this address

	if !k.GetAccountKeeper().HasAccount(ctx, rewardsPoolAccountAddr) {
		fmt.Println("Creating Rewards Pool Module Account (rewards_pool) at address:", rewardsPoolAccountAddr)
		account := k.GetAccountKeeper().NewAccountWithAddress(ctx, rewardsPoolAccountAddr)
		k.GetAccountKeeper().SetAccount(ctx, account)
	} else {
		fmt.Println("Rewards Pool Module Account (rewards_pool) already exists at address:", rewardsPoolAccountAddr)
	}

	// Check if the rewards module account is registered with the bank keeper
	fmt.Println("Module account registered with bank keeper: ", k.GetAccountKeeper().HasAccount(ctx, rewardsModuleAccountAddr))

	// Mint initial EGV tokens to the rewards module account
	mintAmount := math.NewInt(10000)
	err := k.GetBankKeeper().MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(types.EGVDenom, mintAmount)))
	if err != nil {
		panic(fmt.Sprintf("Failed to mint initial EGV tokens: %v", err))
	}
	fmt.Println("Minted", mintAmount, "EGV tokens")

	// Distribute the initial EGV tokens to the predefined addresses based on the genesis state
	predefinedAddresses := genState.Params.PredefinedWallets
	initialEGVAmount := sdk.NewCoin(types.EGVDenom, math.NewInt(1000)) // 1000 EGV tokens to distribute to each predefined address

	for _, addr := range predefinedAddresses {
		address, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			panic(fmt.Sprintf("invalid predefined address: %s", addr))
		}

		// Send initial EGV tokens to each predefined address
		if err := k.SendInitialEGV(ctx, address, initialEGVAmount); err != nil {
			panic(fmt.Sprintf("failed to send initial EGV to address: %s", addr))
		}
	}

	supply := k.GetBankKeeper().GetSupply(ctx, types.EGVDenom)
	fmt.Println("Total supply of EGV:", supply.Amount.String())
}

// ExportGenesis returns the module's exported genesis state as raw JSON bytes.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	// Get the current parameters from the keeper
	params := k.GetParams(ctx)

	// Generate the GenesisState with the current parameters
	genesisState := types.GenesisState{
		Params: params,
	}

	// Return the exported GenesisState
	return &genesisState
}
