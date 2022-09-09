package main

import "fmt"

func FindHomeAction(game *GameStruct) (result []Action) {
	// search free
	for i, c := range game.Free {
		if c > 0 && CanPlaceHome(game, c) {
			a := Action{
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
		if CanPlaceHome(game, card) {
			a := Action{
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
func FindFreeAction(game *GameStruct) (result []Action) {
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
			if GetSequnceLength(game, j) > 4-i {
				continue
			}

			result = append(result, Action{
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
func FindMoveTarget(game *GameStruct, card int, size int) (result []int) {
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
		if CanPlaceOn(card, tarCard) && CanMove(game, size) {
			result = append(result, targetIndex)
		}
	}
	return
}

// 生成移动
func FindMoveAction(game *GameStruct) (result []Action) {
	for index, card := range game.Free {
		movTar := FindMoveTarget(game, card, 1)
		for _, mov := range movTar {
			result = append(result, Action{
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
			if lastCard != 0 && !CanPlaceOn(lastCard, curCard) {
				break
			}
			lastCard = curCard

			movTar := FindMoveTarget(game, curCard, groupLength-cardIndex)
			for _, mov := range movTar {
				result = append(result, Action{
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
func DoAction(game *GameStruct, action *Action) (result GameStruct) {
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
		result.Card[action.TCol] = CombineSlices(result.Card[action.TCol], cards)
	} else if action.Action == "FreeHome" {
		card := result.Free[action.FCol]
		result.Free[action.FCol] = 0
		result.Home[action.TCol] = card
	} else if action.Action == "FreeMove" {
		card := result.Free[action.FCol]
		result.Free[action.FCol] = 0
		result.Card[action.TCol] = CombineSlices(result.Card[action.TCol], []int{card})
	}

	return
}

// 标记访问记录
var Mark map[string]int = make(map[string]int)
var SolverCount int = 0

func DFSSolver(game *GameStruct) []Action {
	/*
		深搜
		1、按Home，Move，Free行动进行搜索
		2、增加全局缓存避免相同场景，并记录结果，记得加锁
		3、找到结果，则返回行动链

		返回倒序
	*/
	// 不重复查找
	sign := HashGame(game)
	if _, ok := Mark[sign]; ok {
		return nil
	}
	Mark[sign] = 1

	SolverCount++
	if SolverCount > 50000 {
		return nil
	}

	TrySolve := func(act []Action) []Action {
		for _, a := range act {
			tmpGame := DoAction(game, &a)

			// last step
			if IsGameFinished(&tmpGame) {
				return []Action{a}
			}

			// search deeper
			tmpResult := DFSSolver(&tmpGame)
			if len(tmpResult) > 0 {
				return append(tmpResult, a)
			}
		}
		return nil
	}

	var act []Action
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
func BestFirstScore(game *GameStruct) int {
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

func BestFirstSolver(game *GameStruct) []Action {
	/*
		维护行动堆
		优先挑选home张多，free张少的进行
	*/
	heap := MinHeap{}
	heap.Add(Node{
		Game:  game,
		Score: 0,
	})
	var cache map[string]int = make(map[string]int)
	var calculation int = 0

	var result *Node
	for !heap.IsEmpty() {
		node := heap.Pop()
		hash := HashGame(node.Game)

		// 该场面计算过，且优于目前场景
		if _, ok := cache[hash]; ok {
			continue
		}
		// fmt.Println(hash, node.Action)
		// PrintGame(node.Game)
		cache[hash] = node.Move

		calculation += 1
		step := node.Move + 1
		var act []Action
		act = append(act, FindHomeAction(node.Game)...)
		act = append(act, FindMoveAction(node.Game)...)
		act = append(act, FindFreeAction(node.Game)...)
		for _, a := range act {
			tmp := DoAction(node.Game, &a)
			n := Node{
				Game:   &tmp,
				Action: CombineActionSlices(node.Action, []Action{a}),
				Score:  -(BestFirstScore(&tmp)*10000 - step),
				Move:   step,
			}
			if IsGameFinished(&tmp) {
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

	fmt.Println("game")
}
