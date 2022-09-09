package main

import (
	"fmt"
	"testing"
	"time"
)

func CreateNewGame() GameStruct {
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

func CreateMiddleGame() GameStruct {
	var game GameStruct = GameStruct{
		Free: [4]int{413, 308, 209, 0},
		Home: [4]int{103, 202, 302, 403},
		Card: [8][]int{
			{109, 311, 203, 304},
			{313, 412, 111, 410},
			{408, 310, 307, 206, 105, 404},
			{210, 309, 208, 107, 406, 305, 204},
			{213, 312, 211},
			{303, 407, 106, 205},
			{112, 411, 110, 409, 108, 207, 306, 405, 104},
			{113, 212},
		},
	}
	return game
}

func CreateLateGame() GameStruct {
	var game GameStruct = GameStruct{
		Home: [4]int{112, 213, 313, 413},
		Card: [8][]int{
			{113},
		},
	}
	return game
}

func TestIsGameFinished(t *testing.T) {
	game := CreateLateGame()
	act := FindHomeAction(&game)
	game = DoAction(&game, &act[0])
	PrintGame(&game)
	fmt.Println(IsGameFinished(&game))
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
		{
			card:   0,
			target: 101,
			result: false,
		},
		{
			card:   101,
			target: 0,
			result: false,
		},
	}

	for _, c := range testCase {
		if CanPlaceOn(c.card, c.target) != c.result {
			t.Error(c)
		}
	}
}

func TestCanPlaceHome(t *testing.T) {
	game := CreateNewGame()
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
	if CanPlaceHome(&game, 0) {
		t.Error("0")
	}
}

func TestFindHomeAction(t *testing.T) {
	game := CreateNewGame()
	game.Card[7] = append(game.Card[7], 101)
	PrintGame(&game)
	result := FindHomeAction(&game)
	for _, r := range result {
		fmt.Println(game.Card[r.FCol][r.FRow], r.Action, r.TCol+1)
	}

	game = CreateMiddleGame()
	PrintGame(&game)
	result = FindHomeAction(&game)
	for _, r := range result {
		fmt.Printf("%8s From %d, %d To %d\n", r.Action, r.FCol, r.FRow, r.TCol)
		fmt.Println(game.Card[r.FCol][r.FRow], r.Action, r.TCol+1)
	}
}

func TestFindFreeAction(t *testing.T) {
	game := CreateNewGame()
	result := FindFreeAction(&game)
	for _, r := range result {
		fmt.Println(game.Card[r.FCol][r.FRow], r.Action, r.TCol)
	}
}

func TestFindMoveAction(t *testing.T) {
	game := CreateNewGame()
	result := FindMoveAction(&game)
	for _, r := range result {
		fmt.Println(game.Card[r.FCol][r.FRow], r.Action, r.TCol)
	}
}

func TestGetSeq(t *testing.T) {
	game := CreateMiddleGame()
	PrintGame(&game)
	fmt.Println(GetSequnceLength(&game, 2))
	fmt.Println(GetSequnceLength(&game, 3))
}

func TestDoAction(t *testing.T) {
	fmt.Println("=============Test Move Card================")
	var game, compareGame GameStruct
	game = CreateNewGame()
	compareGame = game
	result := FindMoveAction(&game)
	rgame := DoAction(&game, &result[0])
	PrintGame(&rgame)
	if !CheckEqual(&game, &compareGame) {
		t.Error("Move Change")
		PrintGame(&game)
		return
	}

	fmt.Println("=============Test Move Home================")
	game = CreateNewGame()
	game.Card[7] = append(game.Card[7], 101)
	compareGame = game
	result = FindHomeAction(&game)
	rgame = DoAction(&game, &result[0])
	PrintGame(&rgame)
	if !CheckEqual(&game, &compareGame) {
		t.Error("Home Change")
		PrintGame(&game)
		return
	}

	fmt.Println("=============Test Move Free================")
	game = CreateNewGame()
	compareGame = game
	result = FindFreeAction(&game)
	rgame = DoAction(&game, &result[0])
	PrintGame(&rgame)
	if !CheckEqual(&game, &compareGame) {
		t.Error("Free Change")
		PrintGame(&game)
		return
	}

	fmt.Println("=============Test Free Home================")
	game = CreateNewGame()
	game.Free[3] = 101
	compareGame = game
	result = FindHomeAction(&game)
	rgame = DoAction(&game, &result[0])
	PrintGame(&rgame)
	if !CheckEqual(&game, &compareGame) {
		t.Error("FreeHome Change")
		PrintGame(&game)
		return
	}

	fmt.Println("=============Test Free Move================")
	game = CreateNewGame()
	game.Free[3] = 307
	compareGame = game
	result = FindMoveAction(&game)
	rgame = DoAction(&game, &result[0])
	PrintGame(&rgame)
	if !CheckEqual(&game, &compareGame) {
		t.Error("FreeMove Change")
		PrintGame(&game)
		return
	}
}

func TestPrintGame(t *testing.T) {
	game := CreateNewGame()
	PrintGame(&game)
}

func TestLateSolver(t *testing.T) {
	game := CreateLateGame()
	PrintGame(&game)
	action := DFSSolver(&game)
	for i, a := range action {
		fmt.Printf("\nStep %03d| %8s From %d, %d To %d\n", i, a.Action, a.FCol, a.FRow, a.TCol)
		game = DoAction(&game, &a)
		PrintGame(&game)
	}
}

func TestMiddleSolver(t *testing.T) {
	game := CreateMiddleGame()
	CheckLegal(&game)
	PrintGame(&game)
	action := DFSSolver(&game)
	fmt.Println(SolverCount)
	for i := len(action) - 1; i >= 0; i-- {
		a := action[i]
		fmt.Printf("Step %03d| %8s From %d, %d To %d\n", i, a.Action, a.FCol, a.FRow, a.TCol)
		game = DoAction(&game, &a)
		// PrintGame(&game)
	}
}

func TestNewGameSolver(t *testing.T) {
	game := CreateNewGame()
	CheckLegal(&game)
	PrintGame(&game)
	action := DFSSolver(&game)
	fmt.Println(SolverCount)
	for i := len(action) - 1; i >= 0; i-- {
		a := action[i]
		count := len(action) - i
		fmt.Printf("Step %03d| %8s From %d, %d To %d\n", count, a.Action, a.FCol, a.FRow, a.TCol)
		game = DoAction(&game, &a)
		PrintGame(&game)
	}
}

func TestSlice(t *testing.T) {
	a := make([]int, 0, 8)
	a = append(a, []int{1, 2, 3, 4, 5, 6}...)
	fmt.Println(a)
	b := a
	a = a[:4]
	fmt.Println(a)
	fmt.Println(b)
}

func TestBFSLateSolver(t *testing.T) {
	game := CreateLateGame()
	PrintGame(&game)
	action := BestFirstSolver(&game)
	for i, a := range action {
		fmt.Printf("\nStep %03d| %8s From %d, %d To %d\n", i, a.Action, a.FCol, a.FRow, a.TCol)
		game = DoAction(&game, &a)
		PrintGame(&game)
	}
}

func TestBFSMiddleSolver(t *testing.T) {
	game := CreateMiddleGame()
	CheckLegal(&game)
	PrintGame(&game)
	action := BestFirstSolver(&game)
	fmt.Println(SolverCount)
	for i, a := range action {
		fmt.Printf("Step %03d| %8s From %d, %d To %d\n", i, a.Action, a.FCol, a.FRow, a.TCol)
		game = DoAction(&game, &a)
		PrintGame(&game)

		time.Sleep(250 * time.Millisecond)
		fmt.Print("\033[H\033[2J")
	}
}

func TestBFSNewGameSolver(t *testing.T) {
	game := CreateNewGame()
	CheckLegal(&game)
	PrintGame(&game)
	action := BestFirstSolver(&game)
	fmt.Println(SolverCount, len(Mark))
	for i, a := range action {
		fmt.Printf("Step %03d| %8s From %d, %d To %d\n", i, a.Action, a.FCol, a.FRow, a.TCol)
		game = DoAction(&game, &a)
		PrintGame(&game)

		// time.Sleep(250 * time.Millisecond)
		// fmt.Print("\033[H\033[2J")
	}
}
