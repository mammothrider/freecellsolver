package solver

import (
	"fmt"
	"freecellsolver/models"
	"freecellsolver/utils"
	"testing"
	"time"
)

func CreateNewGame() models.GameStruct {
	var game models.GameStruct = models.GameStruct{
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

func CreateNewGame2() models.GameStruct {
	var game models.GameStruct = models.GameStruct{
		Card: [8][]int{
			{103, 405, 104, 304, 411, 413, 113},
			{403, 402, 110, 107, 412, 213, 202},
			{109, 205, 306, 310, 408, 311, 101},
			{204, 108, 201, 111, 208, 410, 106},
			{211, 206, 308, 203, 409, 401},
			{112, 309, 407, 307, 302, 105},
			{404, 212, 102, 313, 301, 305},
			{207, 406, 303, 209, 312, 210},
		},
	}
	return game
}

// 非常复杂
func CreateNewGame3() models.GameStruct {
	var game models.GameStruct = models.GameStruct{
		Card: [8][]int{
			{307, 213, 105, 401, 310, 112, 109},
			{411, 304, 308, 412, 406, 206, 207},
			{413, 301, 204, 103, 113, 306, 409},
			{402, 212, 211, 210, 309, 104, 203},
			{302, 101, 303, 201, 410, 404},
			{405, 110, 205, 313, 106, 111},
			{202, 108, 107, 209, 311, 208},
			{312, 403, 305, 102, 407, 408},
		},
	}
	return game
}

// 必须要home操作才能解
func CreateNewGame4() models.GameStruct {
	var game models.GameStruct = models.GameStruct{
		Card: [8][]int{
			{404, 208, 102, 406, 112, 307, 306},
			{206, 304, 403, 408, 412, 104, 303},
			{207, 110, 302, 313, 103, 308, 205},
			{209, 211, 105, 301, 311, 310, 106},
			{407, 111, 410, 305, 107, 405},
			{212, 402, 101, 109, 309, 213},
			{201, 202, 413, 312, 210, 401},
			{113, 108, 204, 409, 411, 203},
		},
	}
	return game

}

func CreateMiddleGame() models.GameStruct {
	var game models.GameStruct = models.GameStruct{
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

func CreateLateGame() models.GameStruct {
	var game models.GameStruct = models.GameStruct{
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
	DoAction(&game, &act[0])
	utils.PrintGame(&game)
	fmt.Println(utils.IsGameFinished(&game))
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
		if utils.CanPlaceOn(c.card, c.target) != c.result {
			t.Error(c)
		}
	}
}

func TestCanPlaceHome(t *testing.T) {
	game := CreateNewGame()
	if !utils.CanPlaceHome(&game, 101) {
		t.Error("101")
	}
	game.Home[0] = 101
	if !utils.CanPlaceHome(&game, 102) {
		t.Error("102")
	}
	if utils.CanPlaceHome(&game, 202) {
		t.Error("202")
	}
	if utils.CanPlaceHome(&game, 0) {
		t.Error("0")
	}
}

func TestFindHomeAction(t *testing.T) {
	game := CreateNewGame()
	game.Card[7] = append(game.Card[7], 101)
	utils.PrintGame(&game)
	result := FindHomeAction(&game)
	for _, r := range result {
		fmt.Println(game.Card[r.FCol][r.FRow], r.Action, r.TCol+1)
	}

	game = CreateMiddleGame()
	utils.PrintGame(&game)
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
	fmt.Println("+++++++++++++++++++++++++++")

	game = models.GameStruct{
		Card: [8][]int{
			{},
			{213, 312},
			{203, 110, 308, 407, 106, 405, 304},
			{102, 301, 307, 411, 313, 105, 413, 112, 211, 310, 209, 108, 207, 306, 205, 104, 403, 302},
			{406, 305, 404, 303, 202},
			{204, 103, 402},
			{409, 212, 111, 210, 109, 408, 107, 206},
			{113, 412, 311, 410, 309, 208},
		},
	}
	result = FindMoveAction(&game)
	for _, r := range result {
		fmt.Println(game.Card[r.FCol][r.FRow], r.Action, r.TCol)
	}
}

func TestGetSeq(t *testing.T) {
	game := CreateMiddleGame()
	utils.PrintGame(&game)
	fmt.Println(utils.GetSequnceLength(&game, 2))
	fmt.Println(utils.GetSequnceLength(&game, 3))
}

func TestDoAction(t *testing.T) {
	fmt.Println("=============Test Move Card================")
	var game, compareGame models.GameStruct
	game = CreateNewGame()
	compareGame = game
	result := FindMoveAction(&game)
	DoAction(&game, &result[0])
	utils.PrintGame(&game)
	if !utils.CheckEqual(&game, &compareGame) {
		t.Error("Move Change")
		utils.PrintGame(&game)
		return
	}

	fmt.Println("=============Test Move Home================")
	game = CreateNewGame()
	game.Card[7] = append(game.Card[7], 101)
	compareGame = game
	result = FindHomeAction(&game)
	DoAction(&game, &result[0])
	utils.PrintGame(&game)
	if !utils.CheckEqual(&game, &compareGame) {
		t.Error("Home Change")
		utils.PrintGame(&game)
		return
	}

	fmt.Println("=============Test Move Free================")
	game = CreateNewGame()
	compareGame = game
	result = FindFreeAction(&game)
	DoAction(&game, &result[0])
	utils.PrintGame(&game)
	if !utils.CheckEqual(&game, &compareGame) {
		t.Error("Free Change")
		utils.PrintGame(&game)
		return
	}

	fmt.Println("=============Test Free Home================")
	game = CreateNewGame()
	game.Free[3] = 101
	compareGame = game
	result = FindHomeAction(&game)
	DoAction(&game, &result[0])
	utils.PrintGame(&game)
	if !utils.CheckEqual(&game, &compareGame) {
		t.Error("FreeHome Change")
		utils.PrintGame(&game)
		return
	}

	fmt.Println("=============Test Free Move================")
	game = CreateNewGame()
	game.Free[3] = 307
	compareGame = game
	result = FindMoveAction(&game)
	DoAction(&game, &result[0])
	utils.PrintGame(&game)
	if !utils.CheckEqual(&game, &compareGame) {
		t.Error("FreeMove Change")
		utils.PrintGame(&game)
		return
	}
}

func TestPrintGame(t *testing.T) {
	game := CreateNewGame()
	utils.PrintGame(&game)
}

func TestLateSolver(t *testing.T) {
	game := CreateLateGame()
	utils.PrintGame(&game)
	action := DFSSolver(&game)
	for i, a := range action {
		fmt.Printf("\nStep %03d| %8s From %d, %d To %d\n", i, a.Action, a.FCol, a.FRow, a.TCol)
		DoAction(&game, &a)
		utils.PrintGame(&game)
	}
}

func TestMiddleSolver(t *testing.T) {
	game := CreateMiddleGame()
	utils.CheckLegal(&game)
	utils.PrintGame(&game)
	action := DFSSolver(&game)
	fmt.Println(SolverCount)
	for i := len(action) - 1; i >= 0; i-- {
		a := action[i]
		fmt.Printf("Step %03d| %8s From %d, %d To %d\n", i, a.Action, a.FCol, a.FRow, a.TCol)
		DoAction(&game, &a)
		// utils.PrintGame(&game)
	}
}

func TestNewGameSolver(t *testing.T) {
	game := CreateNewGame3()
	utils.CheckLegal(&game)
	utils.PrintGame(&game)
	action := DFSSolver(&game)
	fmt.Println(SolverCount)
	for i := len(action) - 1; i >= 0; i-- {
		a := action[i]
		count := len(action) - i
		fmt.Printf("Step %03d| %8s From %d, %d To %d\n", count, a.Action, a.FCol, a.FRow, a.TCol)
		DoAction(&game, &a)
		utils.PrintGame(&game)
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
	utils.PrintGame(&game)
	action := BestFirstSolver(&game)
	for i, a := range action {
		fmt.Printf("\nStep %03d| %8s From %d, %d To %d\n", i, a.Action, a.FCol, a.FRow, a.TCol)
		DoAction(&game, &a)
		utils.PrintGame(&game)
	}
}

func TestBFSMiddleSolver(t *testing.T) {
	game := CreateMiddleGame()
	utils.CheckLegal(&game)
	utils.PrintGame(&game)
	action := BestFirstSolver(&game)
	fmt.Println(SolverCount)
	for i, a := range action {
		fmt.Printf("Step %03d| %8s From %d, %d To %d\n", i, a.Action, a.FCol, a.FRow, a.TCol)
		DoAction(&game, &a)
		utils.PrintGame(&game)

		time.Sleep(250 * time.Millisecond)
		fmt.Print("\033[H\033[2J")
	}
}

func TestBFSNewGameSolver(t *testing.T) {
	game := CreateNewGame2()
	utils.CheckLegal(&game)
	utils.PrintGame(&game)
	action := BestFirstSolver(&game)
	for i, a := range action {
		fmt.Printf("Step %03d| %8s From %d, %d To %d\n", i, a.Action, a.FCol, a.FRow, a.TCol)
		DoAction(&game, &a)

		// utils.PrintGame(&game)
		// time.Sleep(250 * time.Millisecond)
		// fmt.Print("\033[H\033[2J")
	}

}

func TestMultiThreadMiddleGameSolver(t *testing.T) {
	game := CreateMiddleGame()
	utils.CheckLegal(&game)
	utils.PrintGame(&game)
	action := MultiThreadBestFirstSolver(&game)
	fmt.Println(SolverCount)
	for i, a := range action {
		fmt.Printf("Step %03d| %8s From %d, %d To %d\n", i, a.Action, a.FCol, a.FRow, a.TCol)
		DoAction(&game, &a)

		// utils.PrintGame(&game)
		// fmt.Print("\033[H\033[2J")
		// time.Sleep(200 * time.Millisecond)
	}
}

func TestMultiThreadNewGameSolver(t *testing.T) {
	game := CreateNewGame3()
	utils.CheckLegal(&game)
	utils.PrintGame(&game)
	action := MultiThreadBestFirstSolver(&game)
	fmt.Println(SolverCount)
	for i, a := range action {
		fmt.Printf("Step %03d| %8s From %d, %d To %d\n", i, a.Action, a.FCol, a.FRow, a.TCol)
		DoAction(&game, &a)

		// utils.PrintGame(&game)
		// fmt.Print("\033[H\033[2J")
		// time.Sleep(200 * time.Millisecond)
	}
}
