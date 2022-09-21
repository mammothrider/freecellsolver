package solver

import (
	"fmt"
	"freecellsolver/minheap"
	"freecellsolver/models"
	"freecellsolver/utils"
)

// 算分。Home一张10分，Free一张扣1分
func BestFirstScore(game *models.GameStruct) int {
	score := 0
	homeMin := getHomeMin(game.Home)
	for _, c := range game.Home {
		// score += c % 100
		if c%100 <= homeMin+1 {
			score += c % 100
		} else {
			score += homeMin + 1
		}
	}
	// score = score * 10
	// for _, c := range game.Free {
	// 	if c != 0 {
	// 		score--
	// 	}
	// }
	return score
}

func BestFirstSolver(game *models.GameStruct) []models.Action {
	/*
		维护行动堆
		优先挑选home张多，free张少的进行
	*/
	heap := minheap.MinHeap{}
	heap.Add(&models.Node{
		Game:  game,
		Score: 0,
	})
	var cache map[string]*models.Node = make(map[string]*models.Node)
	var calculation int = 0

	var result *models.Node
	for !heap.IsEmpty() && calculation < 700000 {
		node := heap.Pop()
		hash := utils.HashGame(node.Game)
		// 该场面计算过，且优于历史场景
		if v, ok := cache[hash]; ok {
			if v.Move > node.Move {
				v.Parent = node.Parent
			}
			continue
		}
		cache[hash] = node

		calculation += 1
		step := node.Move + 1
		var act []models.Action
		act = append(act, FindUpAction(node.Game)...)
		if len(act) == 0 {
			act = append(act, FindMoveAction(node.Game)...)
			act = append(act, FindFreeAction(node.Game)...)
			act = append(act, FindHomeAction(node.Game)...)
		}

		for _, a := range act {
			copyGame := node.Game.Copy()
			DoAction(copyGame, &a)
			n := models.Node{
				Game:   copyGame,
				Action: a,
				Score:  -(BestFirstScore(copyGame)*10000 - step),
				Move:   step,
				Parent: node,
			}
			if utils.IsGameFinished(copyGame) {
				if result == nil || result.Move > n.Move {
					result = &n
				}
				goto END
			}
			// fmt.Printf("%8s From %d, %d To %d\n", a.Action, a.FCol, a.FRow, a.TCol)
			heap.Add(&n)
		}
	}
END:
	fmt.Println("Total Step:", calculation)
	actions := make([]models.Action, 0, 70)
	for result != nil && result.Parent != nil {
		actions = append(actions, result.Action)
		result = result.Parent
	}
	// reverse
	for i, j := 0, len(actions)-1; i < j; i, j = i+1, j-1 {
		actions[i], actions[j] = actions[j], actions[i]
	}

	return actions
}
