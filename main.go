package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"time"
)

// Merkel Tree
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

// Blockchain
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

// This function calculates the hash representation of the transactions in the block.
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Data) // Assuming the data itself is a hash for simplification.
	}
	txHash := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

// This function sets the hash for a block.
func (b *Block) SetHash() {
	timestamp := []byte(fmt.Sprintf("%x", b.Timestamp))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.HashTransactions(), timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func (b *Block) CalculateMerkelRoot() {
	slice := [][]byte{}
	for _, tx := range b.Transactions {
		slice = append(slice, tx.Data)
	}
	mt := NewMerkleTree(slice)

	b.MerkleRoot = mt.RootNode.Data
}

// Function to add a block to the blockchain.
func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := &Block{time.Now().Unix(), transactions, prevBlock.Hash, nil, nil}
	newBlock.SetHash()
	newBlock.CalculateMerkelRoot()
	bc.blocks = append(bc.blocks, newBlock)
}

func NewGenesisBlock() *Block {
	return &Block{time.Now().Unix(), []*Transaction{{[]byte("Genesis Block")}}, []byte{}, []byte{}, nil}
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

// Interface
const (
	Help           = "help"
	AddBlock       = "add_block"
	AddTransaction = "add_transaction"
	ViewBlockchain = "view_blockchain"
	Exit           = "exit"
)

var tempTransactions []*Transaction

// parseCommand parses command-line arguments
func parseCommand(input string, blockchain *Blockchain) {
	parts := strings.Split(input, " ")
	command := parts[0]
	switch command {
	case ViewBlockchain:
		for _, block := range blockchain.blocks {
			fmt.Printf("Prev. Hash: %x\n", block.PrevBlockHash)
			fmt.Printf("Hash: %x\n", block.Hash)
			fmt.Printf("Timestamp: %d\n", block.Timestamp)
			fmt.Printf("MerkleRoot: %x\n", block.MerkleRoot)
			for _, tx := range block.Transactions {
				fmt.Printf("Transaction Data: %s\n", tx.Data)
			}
			fmt.Println("-------------------------------")
		}

	case AddTransaction:
		tx := &Transaction{[]byte(strings.Join(parts[1:], " "))}
		tempTransactions = append(tempTransactions, tx)
		fmt.Println("Add transaction successfully !")

	case AddBlock:
		if checkempty(tempTransactions) == false {
			fmt.Println("Don`t have transaction!")
		} else {
			blockchain.AddBlock(tempTransactions)
			tempTransactions = nil
			fmt.Println("Add block successfully !")
		}
	case Exit:
		fmt.Println("Exiting...")
		return

	case Help:
		help()

	default:
		fmt.Println("Invalid command. Type 'help' for a list of commands.")
	}
}

func help() {
	fmt.Println("help")
	fmt.Println("add_block")
	fmt.Println("add_transaction 'chuoi muon them'!")
	fmt.Println("view_blockchain")
}

func checkempty(tempTransactions []*Transaction) bool {
	if len(tempTransactions) == 0 {
		return false
	}
	return true
}

func main() {
	help()
	bc := NewBlockchain()
	reader := bufio.NewReader(os.Stdin)
	for {
		// Get user input
		fmt.Print("\nEnter command: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Exiting...")
			break
		}
		parseCommand(input, bc)
	}
}
