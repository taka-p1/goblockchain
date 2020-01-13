package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []string
}

func NewBlock(nonce int, previousHash [32]byte) *Block {
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
	fmt.Printf("previousHash     %x\n", b.previousHash)
	fmt.Printf("transactions     %s\n", b.transactions)
}

type BlockChain struct {
	transactionPool []string
	chain           []*Block
}

func (bc *BlockChain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
	return b
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)         // struct -> json
	return sha256.Sum256([]byte(m)) // json -> struct
}

// MarshalJSON converts Block to json
// To access the fields, use upper case but use lower to convert json.
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64    `json:"timestamp`      // Covert to capital
		Nonce        int      `json:"nonce"`         // Covert to capital
		PreviousHash [32]byte `json:"previous_hash"` // Covert to capital
		Transactions []string `json:"transactions"`  // Covert to capital
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

func NewBlockChain() *BlockChain {
	b := &Block{}
	bc := new(BlockChain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *BlockChain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *BlockChain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s \n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 70))
}

func main() {
	bc := NewBlockChain()
	bc.Print()

	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(5, previousHash)
	bc.Print()
	previousHash = bc.LastBlock().Hash()
	bc.CreateBlock(2, previousHash)
	bc.Print()
}
