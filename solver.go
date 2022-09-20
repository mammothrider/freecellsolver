package main

import (
	"context"
	"encoding/json"
	"fmt"
	"freecellsolver/minheap"
	"freecellsolver/models"
	"freecellsolver/utils"
	"os"
	"runtime"
	"sync"
	"time"
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
	for _, act := range action {
		if game.Home[act.TCol]%100 == min {
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
	heap.Add(&models.Node{
		Game:  game,
		Score: 0,
	})
	var cache map[string]int = make(map[string]int)
	var calculation int = 0

	var result *models.Node
	for !heap.IsEmpty() && calculation < 1000000 {
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
		// act = append(act, FindHomeAction(node.Game)...)
		act = append(act, FindUpAction(node.Game)...)
		act = append(act, FindMoveAction(node.Game)...)
		act = append(act, FindFreeAction(node.Game)...)

		for _, a := range act {
			tmp := DoAction(node.Game, &a)
			n := models.Node{
				Game:   &tmp,
				Action: a,
				Score:  -(BestFirstScore(&tmp)*10000 - step),
				Move:   step,
				Parent: node,
			}
			if utils.IsGameFinished(&tmp) {
				result = &n
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

func ThreadSolve(inputChan, waitChan, resultChan chan *models.Node, cache *sync.Map, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	// defer fmt.Println("ThreadSolve quit")
	for {
		select {
		case <-ctx.Done():
			return
		case node := <-waitChan:
			if node == nil {
				continue
			}
			step := node.Move + 1
			var act []models.Action
			act = append(act, FindUpAction(node.Game)...)
			act = append(act, FindMoveAction(node.Game)...)
			act = append(act, FindFreeAction(node.Game)...)

			tmpResult := make([]*models.Node, 0, len(act))
			for _, a := range act {
				tmp := DoAction(node.Game, &a)
				n := models.Node{
					Game:   &tmp,
					Action: a,
					Score:  -(BestFirstScore(&tmp)*10000 - step),
					Move:   step,
					Parent: node,
				}
				if utils.IsGameFinished(&tmp) {
					resultChan <- &n
					return
				}
				hash := utils.HashGame(n.Game)
				// 该场面计算过
				if _, ok := cache.Load(hash); ok {
					continue
				}
				cache.Store(hash, n.Move)
				tmpResult = append(tmpResult, &n)
			}
			for _, n := range tmpResult {
				select {
				case <-ctx.Done():
					return
				case inputChan <- n:
				}
			}
		}
	}
}

func MultiThreadBestFirstSolver(game *models.GameStruct) []models.Action {
	/*
		inputChan -> heap -> waitChan -> solver -> inputChan
	*/
	var (
		cpuCoreCount int               = runtime.NumCPU()
		inputChan    chan *models.Node = make(chan *models.Node, cpuCoreCount-1)
		waitChan     chan *models.Node = make(chan *models.Node, cpuCoreCount-1)
		resultChan   chan *models.Node = make(chan *models.Node, 1)
		wg           sync.WaitGroup
	)
	ctx, ctxCancel := context.WithCancel(context.Background())
	inputChan <- &models.Node{
		Game:  game,
		Score: 0,
	}

	// 写入到堆里,预处理,并等待分配
	wg.Add(1)
	go func(inputChan, waitChan, resultChan chan *models.Node, ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		// defer fmt.Println("dispatcher quit")
		heap := minheap.MinHeap{}
		calculation := 0
		for {
			timeout := time.After(time.Second)
			if calculation > 1000000 {
				resultChan <- nil
				return
			}
			tmp_node := heap.Get()
			if tmp_node == nil {
				select {
				case <-ctx.Done():
					return
				case node, ok := <-inputChan:
					if !ok {
						return
					}
					heap.Add(node)
				case <-timeout:
					resultChan <- nil
					return
				}
			} else {
				select {
				case <-ctx.Done():
					return
				case node, ok := <-inputChan:
					if !ok {
						return
					}
					heap.Add(node)
				case waitChan <- tmp_node:
					calculation++
					if calculation%10000 == 0 {
						fmt.Println(calculation)
					}
					heap.Pop()
				case <-timeout:
					resultChan <- nil
					return
				}
			}
		}
	}(inputChan, waitChan, resultChan, ctx, &wg)

	// 启动solver
	var cache sync.Map
	for i := 0; i < cpuCoreCount-1; i++ {
		wg.Add(1)
		go ThreadSolve(inputChan, waitChan, resultChan, &cache, ctx, &wg)
	}

	node := <-resultChan
	ctxCancel()

	// empty all channel
EMPTY:
	for {
		select {
		case <-resultChan:
		case <-waitChan:
		case <-inputChan:
		default:
			break EMPTY
		}
	}

	actions := make([]models.Action, 0, 70)
	for node != nil && node.Parent != nil {
		actions = append(actions, node.Action)
		node = node.Parent
	}
	// reverse
	for i, j := 0, len(actions)-1; i < j; i, j = i+1, j-1 {
		actions[i], actions[j] = actions[j], actions[i]
	}
	wg.Wait()
	fmt.Println("Finished")
	close(inputChan)
	close(waitChan)
	close(resultChan)
	return actions
}

func SolveJson(input string) string {
	game := models.GameStruct{}
	err := json.Unmarshal([]byte(input), &game)
	if err != nil {
		panic(err.Error())
	}
	actions := BestFirstSolver(&game)
	res, _ := json.Marshal(actions)
	return string(res)
}

func main() {
	input := os.Args[1]
	fmt.Println(SolveJson(input))
}
