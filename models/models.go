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
