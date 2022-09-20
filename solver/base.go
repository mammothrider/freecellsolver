package solver

import (
	"freecellsolver/models"
	"freecellsolver/utils"
)

// 移动多张到Home
func FindUpAction(game *models.GameStruct) []models.Action {
	action := FindHomeAction(game)
	var tar *models.Action
	min := 99
	for _, h := range game.Home {
		if h%100 < min {
			min = h % 100
		}
	}
	if min == 0 {
		min = 1
	}
	for _, act := range action {
		if game.Home[act.TCol]%100 <= min {
			tar = &act
		}
	}
	if tar != nil {
		tar.Action = "Up"
		return []models.Action{*tar}
	}

	return nil
}

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
			// 只检查一个空位的移动即可
			if toEmpty && utils.CanMove(game, targetIndex, size) {
				toEmpty = false
				result = append(result, targetIndex)
			}
			continue
		}

		tarCard := targetGroup[len(targetGroup)-1]
		if utils.CanPlaceOn(card, tarCard) && utils.CanMove(game, targetIndex, size) {
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
				TRow:   len(game.Card[mov]),
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
					TRow:   len(game.Card[mov]),
				})
			}
		}
	}

	return
}

func DoUpAction(game *models.GameStruct) models.GameStruct {
	copyGame := *game
	action := FindHomeAction(&copyGame)
	do := true
	for action != nil && do {
		min := 99
		for _, h := range copyGame.Home {
			if h%100 < min {
				min = h % 100
			}
		}
		// 2的特殊处理
		if min == 0 {
			min = 1
		}
		do = false
		for _, act := range action {
			if copyGame.Home[act.TCol]%100 <= min {
				copyGame = DoAction(&copyGame, &act)
				do = true
			}
		}

		action = FindHomeAction(&copyGame)
	}
	return copyGame
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
	} else if action.Action == "Up" {
		result = DoUpAction(&result)
	}

	return
}
