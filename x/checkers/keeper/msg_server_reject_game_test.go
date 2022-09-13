package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/smartcoding51/checkers/testutil/keeper"
	"github.com/smartcoding51/checkers/x/checkers"
	"github.com/smartcoding51/checkers/x/checkers/keeper"
	"github.com/smartcoding51/checkers/x/checkers/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServerWithOneGameForRejectGame(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context) {
	k, ctx := keepertest.CheckersKeeper(t)
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)
	server.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Red:     bob,
		Black:   carol,
	})
	return server, *k, context
}

func TestRejectGameByRedOneMoveRemovedGame(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForRejectGame(t)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: carol,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator: bob,
		IdValue: "1",
	})
	nextGame, found := keeper.GetNextGame(sdk.UnwrapSDKContext(context))
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		IdValue: 2,
	}, nextGame)
	_, found = keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "1")
	require.False(t, found)
}
