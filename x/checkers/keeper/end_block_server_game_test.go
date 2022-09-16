package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/smartcoding51/checkers/x/checkers/types"
	"github.com/stretchr/testify/require"
)

func TestForfeitUnplayed(t *testing.T) {
	_, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)

	nextGame, found := keeper.GetNextGame(ctx)
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		IdValue:       2,
		FifoHeadIndex: "-1",
		FifoTailIndex: "-1",
	}, nextGame)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-forfeited",
		Attributes: []sdk.Attribute{
			{Key: "game-index", Value: "1"},
			{Key: "winner", Value: "*"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, event)
}

func TestForfeitPlayedOnce(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: bob,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)

	nextGame, found := keeper.GetNextGame(ctx)
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		IdValue:       2,
		FifoHeadIndex: "-1",
		FifoTailIndex: "-1",
	}, nextGame)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 3)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-forfeited",
		Attributes: []sdk.Attribute{
			{Key: "game-index", Value: "1"},
			{Key: "winner", Value: "*"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, event)
}

func TestForfeitPlayedTwice(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: bob,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: carol,
		IdValue: "1",
		FromX:   0,
		FromY:   5,
		ToX:     1,
		ToY:     4,
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	oldDeadline := types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	game1.Deadline = oldDeadline
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	game1, found = keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Game:        "",
		Turn:        "b",
		Black:       bob,
		Red:         carol,
		MoveCount:   uint64(2),
		BeforeIndex: "-1",
		AfterIndex:  "-1",
		Deadline:    oldDeadline,
		Winner:      "r",
	}, game1)

	nextGame, found := keeper.GetNextGame(ctx)
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		IdValue:       2,
		FifoHeadIndex: "-1",
		FifoTailIndex: "-1",
	}, nextGame)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 3)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-forfeited",
		Attributes: []sdk.Attribute{
			{Key: "game-index", Value: "1"},
			{Key: "winner", Value: "r"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, event)
}
