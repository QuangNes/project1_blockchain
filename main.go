package main

import (
	"bufio"
	"fmt"
	"lab1/blockchain"
	"os"
	"strings"
)

func main() {
	blockchain.HelpCommand()
	bc := blockchain.NewBlockchain()
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
		blockchain.ParseCommand(input, bc)
	}
}
