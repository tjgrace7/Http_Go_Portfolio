package LinkedList

import "fmt"

type Node[T comparable] struct {
	left, right *Node[T]
	data        T
}
type LinkedList[T comparable] struct {
	head *Node[T]
	tail *Node[T]
	size int
}

func (l *LinkedList[T]) insertLeft(data T, node *Node[T]) *Node[T] {
	newNode := &Node[T]{data: data, right: node}
	if node == nil {
		newNode = l.PushFront(data)
		return newNode
	}

	if node.left == nil {
		node.left = newNode
		l.head = newNode
	} else {
		node.left.right = newNode
		newNode.left = node.left
		node.left = newNode
	}
	l.size++
	return newNode
}
func (l *LinkedList[T]) insertRight(data T, node *Node[T]) *Node[T] {
	newNode := &Node[T]{data: data, left: node}
	if node == nil {
		newNode = l.PushBack(data)
		return newNode
	}

	if node.right == nil {
		node.right = newNode
		l.tail = newNode
	} else {
		node.right.left = newNode
		newNode.right = node.right
		node.right = newNode
	}
	l.size++
	return newNode

}
func (l *LinkedList[T]) PushFront(data T) *Node[T] {
	newNode := &Node[T]{data: data}
	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		newNode = l.insertLeft(data, l.head)
		l.head = newNode
	}
	return newNode
}
func (l *LinkedList[T]) PushBack(data T) *Node[T] {
	newNode := &Node[T]{data: data}
	if l.tail == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		newNode = l.insertRight(data, l.tail)
		l.tail = newNode
		fmt.Println("New Tail")
	}
	return newNode
}
func (l *LinkedList[T]) InsertIndex(data T, index int) (*Node[T], error) {
	if l.size <= index {
		return nil, fmt.Errorf("Index Out of Range")
	} else {
		currentNode := l.head
		for i := 0; i < index; i++ {
			currentNode = currentNode.right
		}
		fmt.Println("Current Node Data", currentNode.data)
		newNode := l.insertLeft(data, currentNode)
		return newNode, nil
	}
}
func (l *LinkedList[T]) SearchData(data T) (*Node[T], int, error) {
	currentNode := l.head
	for i := 0; i < l.size; i++ {
		if data == currentNode.data {
			return currentNode, i, nil
		}
		if currentNode.right == nil {
			break
		}
		currentNode = currentNode.right
	}
	return nil, -1, fmt.Errorf("Data not in Linked List")
}
func (l *LinkedList[T]) DeleteData(data T) error {
	deleteNode, index, err := l.SearchData(data)
	if err != nil {
		return fmt.Errorf("Data Not Found, Did Not Delete")
	}

	if deleteNode.left == nil && deleteNode.right == nil {
		l.head = nil
		l.tail = nil
		l.size = 0
	} else if deleteNode.left != nil && deleteNode.right != nil {
		deleteNode.left.right = deleteNode.right
		deleteNode.right.left = deleteNode.left
		l.size--
	} else if deleteNode.left != nil && deleteNode.right == nil {
		deleteNode.left.right = nil
		l.tail = deleteNode.left
		l.size--
	} else {
		deleteNode.right.left = nil
		l.head = deleteNode.right
		l.size--
	}
	fmt.Println("Node deleted at index:", index)
	return nil
}

// If target Index is out of Range, it will push it to the back of the LinkedList instead of erroring
func (l *LinkedList[T]) ReIndex(data T, targetIndex int) {
	_, _, err := l.SearchData(data)
	if err != nil {
		fmt.Println(err)
	}
	l.DeleteData(data)
	if targetIndex >= l.size {
		l.PushBack(data)
	} else {
		l.InsertIndex(data, targetIndex)
	}
}
func (l *LinkedList[T]) Display() {
	currentNode := l.head
	for i := 0; i <= l.size; i++ {
		fmt.Println(currentNode.data)
		currentNode = currentNode.right
	}
}
