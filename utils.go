package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spaolacci/murmur3"
)

// 检查牌能否放在另一张牌上
func CanPlaceOn(card, target int) bool {
	if card == 0 || target == 0 {
		return false
	}
	if card%100+1 != target%100 {
		return false
	}
	if (card/100)%2 == (target/100)%2 {
		return false
	}
	return true
}

// 检查能否放到home区
func CanPlaceHome(game *GameStruct, card int) bool {
	if card < 100 {
		return false
	}
	color := card/100 - 1
	if game.Home[color] == 0 && card%100 == 1 {
		return true
	}
	return game.Home[color]+1 == card
}

// 检查是否有足够的空间移动牌
func CanMove(game *GameStruct, count int) bool {
	free := 0
	empty := 0
	for _, c := range game.Free {
		if c == 0 {
			free++
		}
	}
	for _, g := range game.Card {
		if len(g) == 0 {
			empty++
		}
	}
	return (free+1)*((1+empty)*empty/2+1) >= count
}

func IsGameFinished(game *GameStruct) bool {
	for _, c := range game.Home {
		if c%100 != 13 {
			return false
		}
	}
	return true
}

func PrintGame(game *GameStruct) {
	for _, c := range game.Free {
		fmt.Printf("%3d ", c)
	}
	for _, c := range game.Home {
		fmt.Printf("%3d ", c)
	}
	fmt.Println("\n-------------------------------")
	count := 0
	for {
		flag := true
		for _, g := range game.Card {
			if len(g) > count {
				fmt.Printf("%3d ", g[count])
				flag = false
			} else {
				fmt.Print("--- ")
			}
		}
		fmt.Println()
		count++
		if flag {
			break
		}
	}
}

func murmurUint64(val string) uint64 {
	hasher := murmur3.New64()
	hasher.Write([]byte(val))
	return hasher.Sum64()
}

func HashGame(game *GameStruct) string {
	text, err := json.Marshal(game)
	if err != nil {
		panic("Hash Game Error")
	}

	return strconv.Itoa(int(murmurUint64(string(text))))
}
