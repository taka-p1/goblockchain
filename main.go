package main

import (
	"fmt"
	"log"
	"time"
)

type Block struct {
	nonce        int
	previousHash string
	timestamp    int64
	transactions []string
}

func NewBlock(nonce int, previousHash string) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	return b
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func (b *Block) Print() {
	fmt.Printf("timestamp        %d\n", b.timestamp)
	fmt.Printf("nonce            %d\n", b.nonce)
	fmt.Printf("previousHash     %s\n", b.previousHash)
	fmt.Printf("transactions     %s\n", b.transactions)
}

func main() {
	b := NewBlock(0, "Hash1")
	fmt.Println(b)
	b.Print()
}
