package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/smartcoding51/checkers/x/checkers/types"
	"github.com/stretchr/testify/require"
)

func TestPlayMove2Games1MoveHasSavedFifo(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   carol,
		Red:     alice,
	})

	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: bob,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	nextGame1, found := keeper.GetNextGame(ctx)
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		IdValue:       3,
		FifoHeadIndex: "2",
		FifoTailIndex: "1",
	}, nextGame1)
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Game:        "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "r",
		Black:       bob,
		Red:         carol,
		MoveCount:   uint64(1),
		BeforeIndex: "2",
		AfterIndex:  "-1",
	}, game1)
	game2, found := keeper.GetStoredGame(ctx, "2")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "2",
		Game:        "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "b",
		Black:       carol,
		Red:         alice,
		MoveCount:   uint64(0),
		BeforeIndex: "-1",
		AfterIndex:  "1",
	}, game2)
}

func TestPlayMove2Games2MovesHasSavedFifo(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   carol,
		Red:     alice,
	})
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
		IdValue: "2",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	systemInfo1, found := keeper.GetNextGame(ctx)
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		IdValue:       3,
		FifoHeadIndex: "1",
		FifoTailIndex: "2",
	}, systemInfo1)
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Game:        "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "r",
		Black:       bob,
		Red:         carol,
		MoveCount:   uint64(1),
		BeforeIndex: "-1",
		AfterIndex:  "2",
	}, game1)
	game2, found := keeper.GetStoredGame(ctx, "2")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "2",
		Game:        "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "r",
		Black:       carol,
		Red:         alice,
		MoveCount:   uint64(1),
		BeforeIndex: "1",
		AfterIndex:  "-1",
	}, game2)
}
