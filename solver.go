package main

import "fmt"

// 黑红梅方,黑桃k=113
// 右侧是最外
type GameStruct struct {
	Free [4]int   `json:"free"`
	Home [4]int   `json:"home"`
	Card [8][]int `json:"card"`
}

type Action struct {
	FCol   int
	FRow   int
	Action string
	TCol   int
}

// 检查牌能否放在另一张牌上
func CanPlaceOn(card, target int) bool {
	if card%100+1 != target%100 {
		return false
	}
	if (card/100)%2 == (target/100)%2 {
		return false
	}
	return true
}

// 检查是否有足够的空间移动牌
func CanMove(game GameStruct, count int) bool {
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

func IsGameFinished(game GameStruct) bool {
	for _, c := game.Home {
		if c%100 != 13 {
			return false
		}
	}	
	return true
}

// 生成所有行动
func FindLegalAction(game GameStruct) (result []Action) {
	result = make([]Action, 0)

	free := game.Free
	free_count := -1
	for i, v := range free {
		if v == 0 {
			free_count = i
			break
		}
	}

	home := game.Home
	for i, g := range game.Card {
		if len(g) == 0 {
			continue
		}
		card := g[len(g)-1]
		color := card / 100
		leng := len(g)
		// move to home
		if home[color-1]+1 == card {
			a := Action{
				FCol:   i,
				FRow:   leng - 1,
				Action: "Home",
				TCol:   color - 1,
			}
			result = append(result, a)
		}
		// move to free
		if free_count >= 0 {
			result = append(result, Action{
				FCol:   i,
				FRow:   leng - 1,
				Action: "Free",
				TCol:   free_count,
			})
		}
	}

	// move to another group
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


func main() {

	fmt.Println("game")
}
