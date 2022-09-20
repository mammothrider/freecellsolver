package models

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
	TRow   int
}

type Node struct {
	Game   *GameStruct // 当前场面
	Action Action      // 从上一步到这一步的行动
	Score  int         // 目前分数
	Move   int         // 行动数
	Parent *Node       // 父结点
}

func (g *GameStruct) Copy() *GameStruct {
	game := GameStruct{}
	game.Free = g.Free
	game.Home = g.Home
	for i := 0; i < 8; i++ {
		game.Card[i] = make([]int, len(g.Card[i]))
		copy(game.Card[i], g.Card[i])
	}
	return &game
}
