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
		result.Card[action.TCol] = append(result.Card[action.TCol], cards...)
	} else if action.Action == "FreeHome" {
		card := result.Free[action.FCol]
		result.Free[action.FCol] = 0
		result.Home[action.TCol] = card
	} else if action.Action == "FreeMove" {
		card := result.Free[action.FCol]
		result.Free[action.FCol] = 0
		result.Card[action.TCol] = append(result.Card[action.TCol], card)
	}

	return
}

// 标记访问记录
var Mark map[string]int = make(map[string]int)

func Solver(game *GameStruct) (actionStep []Action) {
	/*
		深搜
		1、按Home，Move，Free行动进行搜索
		2、增加全局缓存避免相同场景，并记录结果，记得加锁
		3、找到结果，则返回行动链
	*/

	// 不重复查找
	sign := HashGame(game)
	if _, ok := Mark[sign]; ok {
		return
	}
	Mark[sign] = 1

	var act []Action
	act = append(act, FindHomeAction(game)...)
	act = append(act, FindMoveAction(game)...)
	act = append(act, FindFreeAction(game)...)

	for _, a := range act {
		tmpGame := DoAction(game, &a)

		// last step
		if IsGameFinished(&tmpGame) {
			actionStep = append(actionStep, a)
			return
		}

		// search deeper
		tmpResult := Solver(&tmpGame)
		if len(tmpResult) > 0 {
			actionStep = append(tmpResult, a)
			return
		}
	}

	return
}

func main() {

	fmt.Println("game")
}
