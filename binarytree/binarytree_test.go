package main

import (
	"fmt"
	"testing"
)

func TestNodeInsert(t *testing.T) {
	var root *Node
	values := []int{47, 12, 89, 3, 65, 91, 28, 56, 74, 19, 82, 6, 39, 100, 15, 63, 51, 8, 97, 34, 72, 45, 21, 88, 60, 4, 77, 29, 95, 11, 58, 42, 84, 17, 69, 33, 99, 25, 53, 7, 80, 38, 66, 22, 93, 49, 14, 71, 36, 62}

	for _, v := range values {
		root = insert(root, v)
	}
	display(root)
	deletednode := delete(root, 56)
	fmt.Println("Search for number returned: ", search(root, 8))
	fmt.Println("Deleted Node is ", deletednode.Value)
	display(root)
}
