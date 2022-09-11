
package types

import (
	// "strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/smartcoding51/checkers/x/checkers/rules"
)

const (
	alice = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d3"
	bob   = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8g"
	carol = "cosmos1e0w5t53nrq7p66fye6c8p0ynyhf6y24l4yuxd7"
)

func GetStoredGame1() *StoredGame {
	return &StoredGame{
		Black:   bob,
		Red:     carol,
		Index:   "1",
		Game:    rules.New().String(),
		Turn:    "b",
	}
}

func TestCanGetAddressBlack(t *testing.T) {
	bobAddress, err1 := sdk.AccAddressFromBech32(bob)
	black, err2 := GetStoredGame1().GetBlackAddress()
	require.Equal(t, bobAddress, black)
	require.Nil(t, err2)
	require.Nil(t, err1)
}

func TestGetAddressWrongBlack(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Black = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8h"
	black, err := storedGame.GetBlackAddress()
	require.Nil(t, black)
	require.EqualError(t,
		err,
		"black address is invalid: cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8h: decoding bech32 failed: checksum failed. Expected xqhc8g, got xqhc8h.")
	require.EqualError(t, storedGame.Validate(), err.Error())
}

func TestCanGetAddressRed(t *testing.T) {
	carolAddress, err1 := sdk.AccAddressFromBech32(carol)
	red, err2 := GetStoredGame1().GetRedAddress()
	require.Equal(t, carolAddress, red)
	require.Nil(t, err1)
	require.Nil(t, err2)
}

func TestGetAddressWrongRed(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Red = "cosmos1e0w5t53nrq7p66fye6c8p0ynyhf6y24l4yuxd8"
	red, err := storedGame.GetRedAddress()
	require.Nil(t, red)
	require.EqualError(t,
		err,
		"red address is invalid: cosmos1e0w5t53nrq7p66fye6c8p0ynyhf6y24l4yuxd8: decoding bech32 failed: checksum failed. Expected 4yuxd7, got 4yuxd8.")
	require.EqualError(t, storedGame.Validate(), err.Error())
}
