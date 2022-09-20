package solver

import (
	"context"
	"fmt"
	"freecellsolver/minheap"
	"freecellsolver/models"
	"freecellsolver/utils"
	"sync"
	"time"
)

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
		cpuCoreCount int               = 8 //runtime.NumCPU()
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
