package rewards_test

import (
	"testing"

	keepertest "zenoda/testutil/keeper"
	"zenoda/testutil/nullify"
	rewards "zenoda/x/rewards/module"
	"zenoda/x/rewards/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.RewardsKeeper(t)
	rewards.InitGenesis(ctx, k, genesisState)
	got := rewards.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
