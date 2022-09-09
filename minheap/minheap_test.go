package minheap

import (
	"fmt"
	"testing"
)

func TestMinHeap(t *testing.T) {
	h := MinHeap{}
	h.Add(Node{Score: 5})
	h.Add(Node{Score: 3})
	h.Add(Node{Score: 1})

	for _, n := range h.node {
		fmt.Println(n.Score)
	}

	h.Add(Node{Score: 2})
	h.Add(Node{Score: 4})
	h.Add(Node{Score: 6})
	h.Add(Node{Score: 8})

	for _, n := range h.node {
		fmt.Println(n.Score)
	}

	var n Node
	n = *h.Pop()
	fmt.Println("Pop", n.Score)
	n = *h.Pop()
	fmt.Println("Pop", n.Score)
	n = *h.Pop()
	fmt.Println("Pop", n.Score)
	h.Add(Node{Score: 7})

	for _, n := range h.node {
		fmt.Println(n.Score)
	}

}
