package main

import (
	"encoding/json"
	"fmt"
	"freecellsolver/models"
	"freecellsolver/solver"
	"os"
)

func SolveJson(input string) string {
	game := models.GameStruct{}
	err := json.Unmarshal([]byte(input), &game)
	if err != nil {
		panic(err.Error())
	}
	actions := solver.BestFirstSolver(&game)
	res, _ := json.Marshal(actions)
	return string(res)
}

func CheckResult(input string) {
	game := models.GameStruct{}
	err := json.Unmarshal([]byte(input), &game)
	if err != nil {
		panic(err.Error())
	}
	actions := solver.BestFirstSolver(&game)
	for i, a := range actions {
		fmt.Printf("Step %03d| %8s From %d, %d To %d\n", i, a.Action, a.FCol, a.FRow, a.TCol)
		game = solver.DoAction(&game, &a)
		// utils.PrintGame(&game)
		// fmt.Print("\033[H\033[2J")
		// time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	input := os.Args[1]
	// fmt.Println(SolveJson(input))
	CheckResult(input)
}
