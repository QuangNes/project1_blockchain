package blockchain

import (
	"fmt"
	"strings"
)

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
func ParseCommand(input string, blockchain *Blockchain) {
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
		if Checkempty(tempTransactions) == false {
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
		HelpCommand()

	default:
		fmt.Println("Invalid command. Type 'help' for a list of commands.")
	}
}

func HelpCommand() {
	fmt.Println("help")
	fmt.Println("add_block")
	fmt.Println("add_transaction 'chuoi muon them'!")
	fmt.Println("view_blockchain")
}

func Checkempty(tempTransactions []*Transaction) bool {
	if len(tempTransactions) == 0 {
		return false
	}
	return true
}
