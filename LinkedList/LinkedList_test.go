package LinkedList

import (
	"fmt"
	"testing"
)

// Int LinkedList Insert

func TestIntegerLinkedList(t *testing.T) {
	values := []int{4, 7, 1, 30, 9, 20, 0, 89, 50}
	link := &LinkedList[int]{}
	for _, v := range values {
		link.PushBack(v)
	}

	node, index, err := link.SearchData(9)
	if err != nil {
		t.Error(err)

	}
	fmt.Println("Node Found ", node, "at index", index)

	node, index, err = link.SearchData(90)
	if err == nil {
		t.Error("Expecting Error Here. Search Data not in list")
	}
	fmt.Println("Node Found ", node, "at index", index)
	err = link.DeleteData(20)
	if err != nil {
		t.Error(err)
	}
	link.PushFront(500)
	link.ReIndex(89, 7)

	link.Display()
}

func TestStringLinkedList(t *testing.T) {
	values := []string{"hello", "world", "tyler", "james", "hannah", "chloe", "programming"}
	link := &LinkedList[string]{}
	for _, v := range values {
		link.PushBack(v)
	}

	node, index, err := link.SearchData("hannah")
	if err != nil {
		t.Error(err)

	}
	fmt.Println("Node Found ", node, "at index", index)

	node, index, err = link.SearchData("dylan")
	if err == nil {
		t.Error("Expecting Error Here. Search Data not in list")
	}
	fmt.Println("Node Found ", node, "at index", index)
	err = link.DeleteData("james")
	if err != nil {
		t.Error(err)
	}

	link.PushFront("Roxy Music")
	link.ReIndex("chloe", 2)
	link.Display()
}
