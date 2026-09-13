package main

import "fmt"

type Node struct {
	Left, Right *Node
	Value       int
}

func insert(node *Node, value int) *Node {

	if node == nil {
		return &Node{Value: value}
	} else if value < node.Value {
		node.Left = insert(node.Left, value)
	} else if value > node.Value {
		node.Right = insert(node.Right, value)
	}
	return node
}

func display(root *Node) {
	if root != nil {
		display(root.Left)
		fmt.Println(root.Value)
		display(root.Right)
	}
}
func search(root *Node, data int) bool {
	if root == nil {
		return false
	} else if root.Value == data {
		return true
	} else if root.Value > data {
		return search(root.Left, data)
	} else {
		return search(root.Right, data)
	}
}
func delete(root *Node, data int) *Node {
	if root == nil {
		return root
	} else if root.Value > data {
		root.Left = delete(root.Left, data)
	} else if root.Value < data {
		root.Right = delete(root.Right, data)
	} else {
		if root.Left == nil && root.Right == nil {
			root = nil
		} else if root.Right != nil { //Find Successor to replace this node
			root.Value = successor(root)
			root.Right = delete(root.Right, root.Value)
		} else {
			root.Value = predecessor(root)
			root.Left = delete(root.Left, root.Value)
		}
	}
	return root
}
func successor(root *Node) int {
	root = root.Right
	for root.Left != nil {
		root = root.Left
	}
	return root.Value
}
func predecessor(root *Node) int {
	root = root.Left
	for root.Right != nil {
		root = root.Right
	}
	return root.Value
}
