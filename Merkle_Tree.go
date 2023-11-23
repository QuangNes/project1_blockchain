package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// MerkleNode represents a node in the Merkle Tree
type MerkleNode struct {
	data  string
	left  *MerkleNode
	right *MerkleNode
}

// NewMerkleNode creates a new MerkleNode with the given data
func NewMerkleNode(data string) *MerkleNode {
	return &MerkleNode{data: data, left: nil, right: nil}
}

// calculateHash computes the hash of the current node based on its data
func (node *MerkleNode) calculateHash() string {
	hash := sha256.Sum256([]byte(node.data))
	return hex.EncodeToString(hash[:])
}

// MerkleTree represents the Merkle Tree structure
type MerkleTree struct {
	root *MerkleNode
}

// NewMerkleTree creates a new MerkleTree with the given data
func NewMerkleTree(data []string) *MerkleTree {
	return &MerkleTree{root: buildTree(data)}
}

// buildTree recursively builds the Merkle Tree
func buildTree(data []string) *MerkleNode {
	if len(data) == 0 {
		return nil
	}

	if len(data) == 1 {
		return NewMerkleNode(data[0])
	}

	var nextLevel []string
	for i := 0; i < len(data); i += 2 {
		left := data[i]
		var right string
		if i+1 < len(data) {
			right = data[i+1]
		}
		combined := left + right
		nextLevel = append(nextLevel, calculateHash(combined))
	}

	return buildTree(nextLevel)
}

// calculateHash computes the SHA-256 hash of the input
func calculateHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// printTree prints the Merkle Tree
func printTree(node *MerkleNode, level int) {
	if node == nil {
		return
	}

	printTree(node.right, level+1)
	fmt.Printf("%*s%s\n", 4*level, "", node.data)
	printTree(node.left, level+1)
}

func main() {
	// Sample data (list of strings)
	data := []string{"Transaction1", "Transaction2", "Transaction3", "Transaction4"}

	// Create a Merkle Tree
	merkleTree := NewMerkleTree(data)

	// Print the Merkle Tree
	fmt.Println("Merkle Tree:")
	printTree(merkleTree.root, 0)
}
