package keeper

import (
	"zenoda/x/rewards/types"

	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DistributeRewards distributes rewards based on transaction volume and predefined governance wallets.
func (k Keeper) DistributeRewards(ctx sdk.Context) {
	// Retrieve parameters
	params := k.GetParams(ctx)

	// Fetch and parse the inflation rate
	inflationRate, err := params.GetInflationRateAsDec()
	if err != nil {
		k.Logger().Error("Invalid inflation rate in params", "error", err)
		return
	}

	// Get total supply of EGV tokens
	totalSupply := k.GetTotalSupply(ctx).Amount

	// Get total network transactions
	totalTx := k.GetTotalTransactions(ctx)
	if totalTx == 0 {
		k.Logger().Info("No transactions recorded on the network; skipping rewards distribution")
		return
	}

	// Iterate over predefined governance addresses and distribute rewards
	for _, walletAddr := range params.PredefinedWallets {
		// Convert wallet address to AccAddress
		addr, err := sdk.AccAddressFromBech32(walletAddr)
		if err != nil {
			k.Logger().Error("Invalid predefined wallet address", "address", walletAddr, "error", err)
			continue
		}

		// Get transaction count for the wallet address
		individualTx := k.GetTransactionCount(ctx, addr)
		if individualTx == 0 {
			k.Logger().Info("No transactions for wallet; skipping reward", "address", addr.String())
			continue
		}

		// Calculate reward using the formula:
		// (individual_address_transactions / total_network_transactions) * (inflation_rate * total_supply)
		reward := inflationRate.MulInt(totalSupply).
			Mul(math.LegacyNewDec(int64(individualTx))).
			Quo(math.LegacyNewDec(int64(totalTx))).
			TruncateInt()

		if reward.IsZero() {
			k.Logger().Info("Calculated reward is zero; skipping distribution", "address", addr.String())
			continue
		}

		// Send reward to the address
		err = k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.ModuleName, addr, sdk.NewCoins(sdk.NewCoin(types.EGVDenom, reward)),
		)
		if err != nil {
			k.Logger().Error("Failed to send reward", "address", addr.String(), "error", err)
			continue
		}

		k.Logger().Info("Reward distributed successfully", "address", addr.String(), "reward", reward.String())
	}
}
