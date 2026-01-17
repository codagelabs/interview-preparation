package main

import "fmt"

type TreeNode struct {
	LeftNode  *TreeNode
	Data      int
	RightNode *TreeNode
}

type Tree struct {
	Root *TreeNode
}

func (tree *Tree) LeftNodeFirstTraverser() {
	currentNode := tree.Root
	stack := []*TreeNode{}
	stack1 := []*TreeNode{}

	stack = append(stack, currentNode)
	for len(stack) > 0 {
		currentNode = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		stack1 = append(stack1, currentNode)

		if currentNode.LeftNode != nil {
			stack = append(stack, currentNode.LeftNode)
		}

		if currentNode.RightNode != nil {
			stack = append(stack, currentNode.RightNode)
		}
	}

	for len(stack1) > 0 {
		currentNode = stack1[len(stack1)-1]
		stack1 = stack1[:len(stack1)-1]
		print(currentNode.Data, " ")
	}
}

//      10
//   2       9
//15   20  11  18

func main() {
	tree := &Tree{
		Root: &TreeNode{
			LeftNode: &TreeNode{
				Data: 2,
				LeftNode: &TreeNode{
					Data: 15,
				},
				RightNode: &TreeNode{
					Data: 20,
				},
			},
			Data: 10,
			RightNode: &TreeNode{
				LeftNode: &TreeNode{
					Data: 11,
				},
				Data: 9,
				RightNode: &TreeNode{
					Data: 18,
				},
			},
		},
	}
	tree.LeftNodeFirstTraverser()
	//InOrderTraversal(tree.Root)
}

func InOrderTraversal(root *TreeNode) {
	stack := []*TreeNode{} // Explicit stack
	current := root

	for current != nil || len(stack) > 0 {
		// Traverse the left subtree
		for current != nil {
			stack = append(stack, current) // Push current node to the stack
			current = current.LeftNode     // Move to the left child
		}

		// Pop the top node from the stack
		current = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Process the current node
		fmt.Print(current.Data, " ")

		// Traverse the right subtree
		current = current.RightNode
	}
}
