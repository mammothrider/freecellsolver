package main

import (
	"fmt"
	"testing"
)

func CreateTestGame() GameStruct {
	var game GameStruct = GameStruct{
		Card: [8][]int{
			{312, 111, 107, 407, 311, 207, 408},
			{304, 211, 302, 309, 208, 106, 406},
			{203, 110, 308, 402, 210, 305, 404},
			{102, 301, 307, 411, 313, 105, 413},
			{101, 103, 306, 205, 104, 202},
			{204, 310, 109, 403, 112, 405},
			{409, 212, 108, 201, 206, 213},
			{113, 401, 412, 303, 209, 410},
		},
	}
	return game
}

func TestCanPlaceOn(t *testing.T) {
	type TestCase struct {
		card, target int
		result       bool
	}

	testCase := []TestCase{
		{
			card:   113,
			target: 112,
			result: false,
		},
		{
			card:   112,
			target: 113,
			result: false,
		},
		{
			card:   102,
			target: 203,
			result: true,
		},
	}

	for _, c := range testCase {
		if CanPlaceOn(c.card, c.target) != c.result {
			t.Error(c)
		}
	}
}

func TestCanPlaceHome(t *testing.T) {
	game := CreateTestGame()
	if !CanPlaceHome(&game, 101) {
		t.Error("101")
	}
	game.Home[0] = 101
	if !CanPlaceHome(&game, 102) {
		t.Error("102")
	}
	if CanPlaceHome(&game, 202) {
		t.Error("202")
	}
}

func TestFindHomeAction(t *testing.T) {
	game := CreateTestGame()
	game.Card[7] = append(game.Card[7], 101)
	result := FindHomeAction(&game)
	for _, r := range result {
		fmt.Println(game.Card[r.FCol][r.FRow], r.Action, r.TCol+1)
	}
}

func TestFindFreeAction(t *testing.T) {
	game := CreateTestGame()
	result := FindFreeAction(&game)
	for _, r := range result {
		fmt.Println(game.Card[r.FCol][r.FRow], r.Action, r.TCol)
	}
}

func TestFindMoveAction(t *testing.T) {
	game := CreateTestGame()
	result := FindMoveAction(&game)
	for _, r := range result {
		fmt.Println(game.Card[r.FCol][r.FRow], r.Action, r.TCol)
	}
}

func TestDoAction(t *testing.T) {
	game := CreateTestGame()
	result := FindMoveAction(&game)
	rgame := DoAction(game, result[0])
	PrintGame(&rgame)

	game.Card[7] = append(game.Card[7], 101)
	result = FindHomeAction(&game)
	rgame = DoAction(game, result[0])
	PrintGame(&rgame)

	result = FindFreeAction(&game)
	rgame = DoAction(game, result[0])
	PrintGame(&rgame)
}

func TestPrintGame(t *testing.T) {
	game := CreateTestGame()
	PrintGame(&game)
}
