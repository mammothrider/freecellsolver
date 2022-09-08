package main

import "fmt"

func FindHomeAction(game *GameStruct) (result []Action) {
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

// 生成移动
func FindMoveAction(game *GameStruct) (result []Action) {
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

			// 检查能不能移动到其他组
			for targetIndex, targetGroup := range game.Card {
				if targetIndex == groupIndex {
					continue
				}
				isEmpty := len(targetGroup) == 0
				tarCard := 0
				if !isEmpty {
					tarCard = targetGroup[len(targetGroup)-1]
				}

				if isEmpty || (CanPlaceOn(curCard, tarCard) && CanMove(game, groupLength-cardIndex)) {
					result = append(result, Action{
						FCol:   groupIndex,
						FRow:   cardIndex,
						Action: "Move",
						TCol:   targetIndex,
					})
				}

			}
		}
	}

	return
}

// 执行动作，返回一个新object
func DoAction(game GameStruct, action Action) GameStruct {
	if action.Action == "Free" {
		leng := len(game.Card[action.FCol])
		card := game.Card[action.FCol][leng-1]
		game.Card[action.FCol] = game.Card[action.FCol][:leng-1]
		game.Free[action.TCol] = card
	} else if action.Action == "Home" {
		leng := len(game.Card[action.FCol])
		card := game.Card[action.FCol][leng-1]
		game.Card[action.FCol] = game.Card[action.FCol][:leng-1]
		game.Home[action.TCol] = card
	} else if action.Action == "Move" {
		cards := game.Card[action.FCol][action.FRow:]
		game.Card[action.FCol] = game.Card[action.FCol][:action.FRow]
		game.Card[action.TCol] = append(game.Card[action.TCol], cards...)
	}

	return game
}

func Solver(game *GameStruct) {
	/*
		深搜
		1、按Home，Move，Free行动进行搜索
		2、增加全局缓存避免相同场景，并记录结果，记得加锁
		3、找到结果，则返回行动链
	*/
}

func main() {

	fmt.Println("game")
}
