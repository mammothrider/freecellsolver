package solver

import (
	"context"
	"fmt"
	"freecellsolver/minheap"
	"freecellsolver/models"
	"freecellsolver/utils"
	"runtime"
	"sync"
	"sync/atomic"
)

func ThreadSolve(heap *minheap.SafeMinHeap, resultChan chan *models.Node, cache *sync.Map, ctx context.Context, wg *sync.WaitGroup, calculation *int32) {
	defer wg.Done()
	// defer fmt.Println("ThreadSolve quit")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			node := heap.Pop()
			if node == nil {
				continue
			}
			atomic.AddInt32(calculation, 1)
			step := node.Move + 1
			var act []models.Action
			act = append(act, FindUpAction(node.Game)...)
			// Up是必须行为
			if len(act) == 0 {
				act = append(act, FindMoveAction(node.Game)...)
				act = append(act, FindFreeAction(node.Game)...)
			}
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
				heap.Add(&n)
			}
		}
	}
}

func MultiThreadBestFirstSolver(game *models.GameStruct) []models.Action {
	/*
		inputChan -> heap -> waitChan -> solver -> inputChan
	*/
	var (
		cpuCoreCount int               = runtime.NumCPU() - 1
		resultChan   chan *models.Node = make(chan *models.Node, 1)
		wg           sync.WaitGroup
		heap         minheap.SafeMinHeap = minheap.SafeMinHeap{}
		calculation  int32
	)
	ctx, ctxCancel := context.WithCancel(context.Background())
	heap.Add(&models.Node{
		Game:  game,
		Score: 0,
	})

	// 启动solver
	var cache sync.Map
	for i := 0; i < cpuCoreCount; i++ {
		wg.Add(1)
		go ThreadSolve(&heap, resultChan, &cache, ctx, &wg, &calculation)
	}

	node := <-resultChan
	ctxCancel()

	// empty all channel
EMPTY:
	for {
		select {
		case <-resultChan:
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
	fmt.Println("Finished, calculate", calculation)
	close(resultChan)
	return actions
}
