package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
	}

	node.Left = left
	node.Right = right

	return &node
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	for _, dat := range data {
		node := NewMerkleNode(nil, nil, dat)
		nodes = append(nodes, *node)
	}

	for i := 0; i < len(data)/2; i++ {
		var level []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			level = append(level, *node)
		}

		nodes = level
	}

	tree := MerkleTree{&nodes[0]}

	return &tree
}

// printTree prints the Merkle Tree
func printTree(node *MerkleNode, level int) {
	if node == nil {
		return
	}
	printTree(node.Right, level+1)
	fmt.Printf("\n", 4*level, "", node.Data)
	printTree(node.Left, level+1)
}

type Transaction struct {
	Data []byte
}

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	MerkleRoot    []byte
}

type Blockchain struct {
	blocks []*Block
}

// This function sets the hash for a block.
func (b *Block) SetHash() {
	timestamp := []byte(fmt.Sprintf("%x", b.Timestamp))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.HashTransactions(), timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

// This function calculates the hash representation of the transactions in the block.
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Data) // Assuming the data itself is a hash for simplification.
	}
	txHash := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

// Function to add a block to the blockchain.
func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := &Block{time.Now().Unix(), transactions, prevBlock.Hash, []byte{}, nil}
	newBlock.SetHash()
	slice := [][]byte{}

	for _, t := range transactions {
		slice = append(slice, t.Data)
	}
	mt := NewMerkleTree(slice)
	newBlock.MerkleRoot = mt.RootNode.Data
	bc.blocks = append(bc.blocks, newBlock)
}

func NewGenesisBlock() *Block {
	return &Block{time.Now().Unix(), []*Transaction{{[]byte("Genesis Block")}}, []byte{}, []byte{}, nil}
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func main() {
	//var data[][]
	bc := NewBlockchain()

	// Add a new block with a single transaction
	tx1 := &Transaction{[]byte("This is transaction 1")}
	bc.AddBlock([]*Transaction{tx1})

	// Add another block with two transactions
	tx2 := &Transaction{[]byte("This is transaction 2")}
	bc.AddBlock([]*Transaction{tx2})

	tx3 := &Transaction{[]byte("This is transaction 3")}
	bc.AddBlock([]*Transaction{tx3})
	// Print details of the blockchain
	for _, block := range bc.blocks {
		fmt.Printf("Prev. Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("MerkleRoot: %x\n", block.MerkleRoot)
		for _, tx := range block.Transactions {
			fmt.Printf("Transaction Data: %s\n", tx.Data)
		}
		fmt.Println("-------------------------------")
	}

}
