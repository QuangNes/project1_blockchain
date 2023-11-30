package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

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
