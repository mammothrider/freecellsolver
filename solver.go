package main

import (
	"fmt"
	"freecellsolver/minheap"
	"freecellsolver/models"
	"freecellsolver/utils"
	"time"
)

func FindHomeAction(game *models.GameStruct) (result []models.Action) {
	// search free
	for i, c := range game.Free {
		if c > 0 && utils.CanPlaceHome(game, c) {
			a := models.Action{
				FCol:   i,
				FRow:   0,
				Action: "FreeHome",
				TCol:   c/100 - 1,
			}
			result = append(result, a)
		}
	}

	// search cards
	for i, g := range game.Card {
		leng := len(g)
		if leng == 0 {
			continue
		}
		card := g[leng-1]
		// move to home
		if utils.CanPlaceHome(game, card) {
			a := models.Action{
				FCol:   i,
				FRow:   leng - 1,
				Action: "Home",
				TCol:   card/100 - 1,
			}
			result = append(result, a)
		}
	}
	return
}

// 移动到Free区
func FindFreeAction(game *models.GameStruct) (result []models.Action) {
	free := game.Free
	for i, v := range free {
		if v != 0 {
			continue
		}

		for j, g := range game.Card {
			leng := len(g)
			if len(g) == 0 {
				continue
			}

			// 序列长度大于空白，禁止移动到Free区
			if utils.GetSequnceLength(game, j) > 4-i {
				continue
			}

			result = append(result, models.Action{
				FCol:   j,
				FRow:   leng - 1,
				Action: "Free",
				TCol:   i,
			})
		}

		break
	}
	return
}

// 检查能不能移动到其他组
func FindMoveTarget(game *models.GameStruct, card int, size int) (result []int) {
	if card == 0 || size == 0 {
		return
	}
	toEmpty := true
	for targetIndex, targetGroup := range game.Card {
		if len(targetGroup) == 0 {
			if toEmpty {
				toEmpty = false
				result = append(result, targetIndex)
			}
			continue
		}

		tarCard := targetGroup[len(targetGroup)-1]
		if utils.CanPlaceOn(card, tarCard) && utils.CanMove(game, size) {
			result = append(result, targetIndex)
		}
	}
	return
}

// 生成移动
func FindMoveAction(game *models.GameStruct) (result []models.Action) {
	for index, card := range game.Free {
		movTar := FindMoveTarget(game, card, 1)
		for _, mov := range movTar {
			result = append(result, models.Action{
				FCol:   index,
				FRow:   0,
				Action: "FreeMove",
				TCol:   mov,
			})
		}
	}

	for groupIndex, curGroup := range game.Card {
		groupLength := len(curGroup)
		if groupLength == 0 {
			continue
		}

		lastCard := 0
		for cardIndex := groupLength - 1; cardIndex >= 0; cardIndex-- {
			curCard := curGroup[cardIndex]

			// 检查是不是在序列里
			if lastCard != 0 && !utils.CanPlaceOn(lastCard, curCard) {
				break
			}
			lastCard = curCard

			movTar := FindMoveTarget(game, curCard, groupLength-cardIndex)
			for _, mov := range movTar {
				result = append(result, models.Action{
					FCol:   groupIndex,
					FRow:   cardIndex,
					Action: "Move",
					TCol:   mov,
				})
			}
		}
	}

	return
}

// 执行动作，返回一个新object
func DoAction(game *models.GameStruct, action *models.Action) (result models.GameStruct) {
	result = *game
	if action.Action == "Free" {
		leng := len(result.Card[action.FCol])
		card := result.Card[action.FCol][leng-1]
		result.Card[action.FCol] = result.Card[action.FCol][:leng-1]
		result.Free[action.TCol] = card
	} else if action.Action == "Home" {
		leng := len(result.Card[action.FCol])
		card := result.Card[action.FCol][leng-1]
		result.Card[action.FCol] = result.Card[action.FCol][:leng-1]
		result.Home[action.TCol] = card
	} else if action.Action == "Move" {
		cards := result.Card[action.FCol][action.FRow:]
		result.Card[action.FCol] = result.Card[action.FCol][:action.FRow]
		result.Card[action.TCol] = utils.CombineSlices(result.Card[action.TCol], cards)
	} else if action.Action == "FreeHome" {
		card := result.Free[action.FCol]
		result.Free[action.FCol] = 0
		result.Home[action.TCol] = card
	} else if action.Action == "FreeMove" {
		card := result.Free[action.FCol]
		result.Free[action.FCol] = 0
		result.Card[action.TCol] = utils.CombineSlices(result.Card[action.TCol], []int{card})
	}

	return
}

// 标记访问记录
var Mark map[string]int = make(map[string]int)
var SolverCount int = 0

func DFSSolver(game *models.GameStruct) []models.Action {
	/*
		深搜
		1、按Home，Move，Free行动进行搜索
		2、增加全局缓存避免相同场景，并记录结果，记得加锁
		3、找到结果，则返回行动链

		返回倒序
	*/
	// 不重复查找
	sign := utils.HashGame(game)
	if _, ok := Mark[sign]; ok {
		return nil
	}
	Mark[sign] = 1

	SolverCount++
	if SolverCount > 50000 {
		return nil
	}

	TrySolve := func(act []models.Action) []models.Action {
		for _, a := range act {
			tmpGame := DoAction(game, &a)

			// last step
			if utils.IsGameFinished(&tmpGame) {
				return []models.Action{a}
			}

			// search deeper
			tmpResult := DFSSolver(&tmpGame)
			if len(tmpResult) > 0 {
				return append(tmpResult, a)
			}
		}
		return nil
	}

	var act []models.Action
	if act = TrySolve(FindHomeAction(game)); act != nil {
		return act
	}
	if act = TrySolve(FindMoveAction(game)); act != nil {
		return act
	}
	if act = TrySolve(FindFreeAction(game)); act != nil {
		return act
	}

	return nil
}

// 算分。Home一张10分，Free一张扣1分
func BestFirstScore(game *models.GameStruct) int {
	score := 0
	for _, c := range game.Home {
		score += c % 100
	}
	score = score * 10
	for _, c := range game.Free {
		if c != 0 {
			score--
		}
	}
	return score
}

func BestFirstSolver(game *models.GameStruct) []models.Action {
	/*
		维护行动堆
		优先挑选home张多，free张少的进行
	*/
	heap := minheap.MinHeap{}
	heap.Add(minheap.Node{
		Game:  game,
		Score: 0,
	})
	var cache map[string]int = make(map[string]int)
	var calculation int = 0

	var result *minheap.Node
	for !heap.IsEmpty() && calculation < 100000 {
		node := heap.Pop()
		hash := utils.HashGame(node.Game)
		// 该场面计算过，且优于目前场景
		if _, ok := cache[hash]; ok {
			continue
		}
		cache[hash] = node.Move

		calculation += 1
		step := node.Move + 1
		var act []models.Action
		act = append(act, FindHomeAction(node.Game)...)
		act = append(act, FindMoveAction(node.Game)...)
		act = append(act, FindFreeAction(node.Game)...)
		for _, a := range act {
			tmp := DoAction(node.Game, &a)
			n := minheap.Node{
				Game:   &tmp,
				Action: utils.CombineActionSlices(node.Action, []models.Action{a}),
				Score:  -(BestFirstScore(&tmp)*10000 - step),
				Move:   step,
			}
			if utils.IsGameFinished(&tmp) {
				result = &n
				goto END
			}
			// fmt.Printf("%8s From %d, %d To %d\n", a.Action, a.FCol, a.FRow, a.TCol)
			heap.Add(n)
		}
	}
END:
	fmt.Println("Total Step:", calculation)
	if result != nil {
		fmt.Println("Move Step:", result.Move)
		return result.Action
	}

	return nil
}

func main() {
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
	utils.CheckLegal(&game)
	utils.PrintGame(&game)
	action := BestFirstSolver(&game)
	for i, a := range action {
		fmt.Printf("Step %03d| %8s From %d, %d To %d\n", i, a.Action, a.FCol, a.FRow, a.TCol)
		game = DoAction(&game, &a)
		utils.PrintGame(&game)

		time.Sleep(250 * time.Millisecond)
		fmt.Print("\033[H\033[2J")
	}
}
