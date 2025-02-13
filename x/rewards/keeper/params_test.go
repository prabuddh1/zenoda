package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "zenoda/testutil/keeper"
	"zenoda/x/rewards/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.RewardsKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
