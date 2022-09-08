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

func GetSequnceLength(game *GameStruct, rol int) int {
	group := game.Card[rol]
	count := 0
	last := 0
	for i := len(group) - 1; i >= 0; i-- {
		if last == 0 || CanPlaceOn(last, group[i]) {
			last = group[i]
			count++
		} else {
			break
		}
	}
	return count
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
	cgame := *game
	for i := 0; i < 4; i++ {
		for j := i; j < 4; j++ {
			if cgame.Free[i] < cgame.Free[j] {
				cgame.Free[i], cgame.Free[j] = cgame.Free[j], cgame.Free[i]
			}
		}
	}
	for i := 0; i < 8; i++ {
		for j := i; j < 8; j++ {
			if len(cgame.Card[i]) < len(cgame.Card[j]) {
				cgame.Card[i], cgame.Card[j] = cgame.Card[j], cgame.Card[i]
			}
		}
	}

	text, err := json.Marshal(cgame)
	if err != nil {
		panic("Hash Game Error")
	}

	return strconv.Itoa(int(murmurUint64(string(text))))
}

func CheckCard(game *GameStruct, card int) bool {
	count := 0
	if card <= game.Home[card/100-1] {
		count++
	}
	for _, c := range game.Free {
		if c == card {
			count++
		}
	}
	for _, g := range game.Card {
		for _, c := range g {
			if c == card {
				count++
			}
		}
	}
	return count == 1
}

func CheckLegal(game *GameStruct) {
	for j := 1; j < 5; j++ {
		for i := 1; i < 14; i++ {
			card := 100*j + i
			if !CheckCard(game, card) {
				err := fmt.Errorf("error %d", card)
				panic(err)
			}
		}
	}
}

func CheckEqual(game1, game2 *GameStruct) bool {
	if game1.Free != game2.Free {
		return false
	}
	if game1.Home != game2.Home {
		return false
	}
	for i, g := range game1.Card {
		if len(game1.Card[i]) != len(game2.Card[i]) {
			return false
		}
		for j, _ := range g {
			if game1.Card[i][j] != game2.Card[i][j] {
				return false
			}
		}
	}
	return true
}

func CombineSlices(slices ...[]int) []int {
	count := 0
	for _, s := range slices {
		count += len(s)
	}
	a := make([]int, 0, count)
	for _, s := range slices {
		a = append(a, s...)
	}
	return a
}
