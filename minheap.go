package main

type Node struct {
	Game   *GameStruct // 当前场面
	Action []Action    // 之前的行动
	Score  int         // 目前分数
	Move   int         // 行动数
}

type MinHeap struct {
	node []Node
}

func (m *MinHeap) Sort() {
	for i := len(m.node) - 1; i >= 0; i-- {
		root := i / 2
		if m.node[root].Score > m.node[i].Score {
			m.node[i], m.node[root] = m.node[root], m.node[i]
		}
	}
}

func (m *MinHeap) Add(n Node) {
	m.node = append(m.node, n)
	for i := len(m.node) - 1; i > 0; i = i / 2 {
		root := i / 2
		if m.node[root].Score > m.node[i].Score {
			m.node[i], m.node[root] = m.node[root], m.node[i]
		}
	}
}

func (m *MinHeap) Pop() *Node {
	if m.node == nil {
		return nil
	}
	root := m.node[0]
	m.node[0] = m.node[len(m.node)-1]
	m.node = m.node[:len(m.node)-1]
	for i := 0; i < len(m.node); {
		l := i * 2
		r := i*2 + 1
		t := i
		if l < len(m.node) && m.node[l].Score < m.node[t].Score {
			t = l
		}
		if r < len(m.node) && m.node[r].Score < m.node[t].Score {
			t = r
		}
		if t == i {
			break
		}

		m.node[i], m.node[t] = m.node[t], m.node[i]
		i = t
	}
	return &root
}

func (m *MinHeap) IsEmpty() bool {
	return len(m.node) == 0
}