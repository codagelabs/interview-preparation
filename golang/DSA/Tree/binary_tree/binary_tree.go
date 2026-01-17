package main

import (
	"fmt"
)

type BinaryTree struct {
	Root *TreeNode
}

type TreeNode struct {
	LeftNode  *TreeNode
	Data      int
	RightNode *TreeNode
}

func Insert(root *TreeNode, data int) *TreeNode {
	if root == nil {
		return &TreeNode{Data: data}
	}
	if data < root.Data {
		root.LeftNode = Insert(root.LeftNode, data)
	} else {
		root.RightNode = Insert(root.RightNode, data)
	}
	return root
}

// InOrderTraversal : Left → Root → Right
func InOrderTraversal(root *TreeNode) {
	if root != nil {
		InOrderTraversal(root.LeftNode)
		fmt.Println(" ", root.Data)
		InOrderTraversal(root.RightNode)
	}
}

// PreOrderTraversal : Root → left → Right
func PreOrderTraversal(root *TreeNode) {
	if root != nil {
		fmt.Println(" ", root.Data)
		PreOrderTraversal(root.LeftNode)
		PreOrderTraversal(root.RightNode)
	}
}

// PostOrderTraversal : Left → Right → Root
func PostOrderTraversal(root *TreeNode) {
	if root != nil {
		PostOrderTraversal(root.LeftNode)
		PostOrderTraversal(root.RightNode)
		fmt.Println(" ", root.Data)
	}
}
func Mirror(root *TreeNode) {
	if root != nil {
		root.LeftNode, root.RightNode = root.RightNode, root.LeftNode
		Mirror(root.LeftNode)
		Mirror(root.RightNode)
	}
}

func CheckIfTwoTreesAreIdentical(root1, root2 *TreeNode) bool {
	if root1 == nil && root2 == nil {
		return true
	}
	if root1 != nil && root2 != nil && root1.Data == root2.Data {
		return CheckIfTwoTreesAreIdentical(root1.LeftNode, root2.LeftNode) && CheckIfTwoTreesAreIdentical(root1.RightNode, root2.RightNode)
	}
	return false
}

func GetHeight(root *TreeNode, searchKey int) int {
	if root == nil {
		return -1
	}

	if root.Data == searchKey {
		return 0
	}
	if searchKey < root.Data {
		return GetHeight(root.LeftNode, searchKey) + 1
	}

	return GetHeight(root.RightNode, searchKey) + 1

}

func GetTreeHeight(root *TreeNode) int {
	if root == nil {
		return -1
	}
	leftHieght := GetTreeHeight(root.LeftNode)
	rightHeight := GetTreeHeight(root.RightNode)
	if leftHieght > rightHeight {
		return leftHieght + 1
	}
	return rightHeight + 1

}

func Search(root *TreeNode, searchKey int) *TreeNode {
	if root == nil || root.Data == searchKey {
		return root
	}

	if searchKey < root.Data {
		return Search(root.LeftNode, searchKey)
	}
	return Search(root.RightNode, searchKey)

}

func FindOutNodeCount(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return 1 + FindOutNodeCount(root.LeftNode) + FindOutNodeCount(root.RightNode)
}

func main() {

	Root := &TreeNode{}
	Insert(Root, 30)
	Insert(Root, 20)
	Insert(Root, 3)
	Insert(Root, 2)
	Insert(Root, 4)
	Insert(Root, 5)
	Insert(Root, 6)
	Insert(Root, 7)
	PostOrderTraversal(Root)
	fmt.Println("search", Search(Root, 2))
	fmt.Println("Height", GetHeight(Root, 7))
	fmt.Println("Tree Height", GetTreeHeight(Root))
	fmt.Println("Total Node", FindOutNodeCount(Root))
	fmt.Println("Check IfT wo Trees Are Identical", CheckIfTwoTreesAreIdentical(Root, Root.LeftNode))

}
