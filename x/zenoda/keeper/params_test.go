package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "zenoda/testutil/keeper"
	"zenoda/x/zenoda/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.ZenodaKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
