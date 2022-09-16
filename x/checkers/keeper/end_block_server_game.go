package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	rules "github.com/smartcoding51/checkers/x/checkers/rules"
	"github.com/smartcoding51/checkers/x/checkers/types"
)

func (k Keeper) ForfeitExpiredGames(goCtx context.Context) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	opponents := map[string]string{
		rules.PieceStrings[rules.BLACK_PLAYER]: rules.PieceStrings[rules.RED_PLAYER],
		rules.PieceStrings[rules.RED_PLAYER]:   rules.PieceStrings[rules.BLACK_PLAYER],
	}

	// Get FIFO information
	nextGame, found := k.GetNextGame(ctx)
	if !found {
		panic("NextGame not found")
	}

	gameIndex := nextGame.FifoHeadIndex
	var storedGame types.StoredGame

	for {
		// Finished moving along
		if gameIndex == types.NoFifoIndex {
			break
		}
		storedGame, found = k.GetStoredGame(ctx, gameIndex)
		if !found {
			panic("Fifo head game not found " + nextGame.FifoHeadIndex)
		}
		deadline, err := storedGame.GetDeadlineAsTime()
		if err != nil {
			panic(err)
		}
		if deadline.Before(ctx.BlockTime()) {
			// Game is past deadline
			k.RemoveFromFifo(ctx, &storedGame, &nextGame)
			lastBoard := storedGame.Game
			if storedGame.MoveCount <= 1 {
				// No point in keeping a game that was never really played
				k.RemoveStoredGame(ctx, gameIndex)
			} else {
				storedGame.Winner, found = opponents[storedGame.Turn]
				if !found {
					panic(fmt.Sprintf(types.ErrCannotFindWinnerByColor.Error(), storedGame.Turn))
				}
				storedGame.Game = ""
				k.SetStoredGame(ctx, storedGame)
			}
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(types.GameForfeitedEventType,
					sdk.NewAttribute(types.GameForfeitedEventGameIndex, gameIndex),
					sdk.NewAttribute(types.GameForfeitedEventWinner, storedGame.Winner),
					sdk.NewAttribute(types.GameForfeitedEventBoard, lastBoard),
				),
			)
			// Move along FIFO
			gameIndex = nextGame.FifoHeadIndex
		} else {
			// All other games after are active anyway
			break
		}
	}
	k.SetNextGame(ctx, nextGame)
}
