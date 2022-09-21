package minheap

import (
	"fmt"
	"freecellsolver/models"
	"testing"
)

func TestMinHeap(t *testing.T) {
	h := MinHeap{}
	h.Add(&models.Node{Score: 5})
	h.Add(&models.Node{Score: 3})
	h.Add(&models.Node{Score: 1})

	for _, n := range h.node {
		fmt.Println(n.Score)
	}

	h.Add(&models.Node{Score: 2})
	h.Add(&models.Node{Score: 4})
	h.Add(&models.Node{Score: 6})
	h.Add(&models.Node{Score: 8})

	for _, n := range h.node {
		fmt.Println(n.Score)
	}

	var n models.Node
	n = *h.Pop()
	fmt.Println("Pop", n.Score)
	n = *h.Pop()
	fmt.Println("Pop", n.Score)
	n = *h.Pop()
	fmt.Println("Pop", n.Score)
	h.Add(&models.Node{Score: 7})

	for _, n := range h.node {
		fmt.Println(n.Score)
	}

}
