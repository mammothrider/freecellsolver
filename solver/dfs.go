package solver

import (
	"freecellsolver/models"
	"freecellsolver/utils"
)

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
	// if act = TrySolve(FindHomeAction(game)); act != nil {
	// 	return act
	// }
	if act = TrySolve(FindUpAction(game)); act != nil {
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
