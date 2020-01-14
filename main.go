package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN" // Adress
	MINING_REWARD     = 1.0
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func (b *Block) Print() {
	fmt.Printf("timestamp        %d\n", b.timestamp)
	fmt.Printf("nonce            %d\n", b.nonce)
	fmt.Printf("previousHash     %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
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
		Timestamp    int64          `json:"timestamp`      // Covert to capital
		Nonce        int            `json:"nonce"`         // Covert to capital
		PreviousHash [32]byte       `json:"previous_hash"` // Covert to capital
		Transactions []*Transaction `json:"transactions"`  // Covert to capital
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

func NewBlockchain(blockchainAddress string) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s \n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 70))
}

type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	t := new(Transaction)
	t.senderBlockchainAddress = sender
	t.recipientBlockchainAddress = recipient
	t.value = value
	return t
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("sender_blockchain_address                %s\n", t.senderBlockchainAddress)
	fmt.Printf("recipient_blockchain_address             %s\n", t.recipientBlockchainAddress)
	fmt.Printf("value                                    %.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`    // Covert to capital
		Recipient string  `json:"recipient_blockchain_address"` // Covert to capital
		Value     float32 `json:"value"`                        // Covert to capital
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(t.senderBlockchainAddress,
				t.recipientBlockchainAddress,
				t.value))
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	// Repeat until ValidProof returns true(find answer)
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) (totalAmount float32) {
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipientBlockchainAddress {
				totalAmount += value
			}
			if blockchainAddress == t.senderBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

func main() {
	myBlockChainAddress := "my_blockchain_address" // In this case, me = miner.
	bc := NewBlockchain(myBlockChainAddress)
	bc.Print()

	bc.AddTransaction("Mr. A", "B-san", 1.0)
	bc.Mining()
	bc.Print()

	bc.AddTransaction("C", "D", 2.0)
	bc.AddTransaction("X", "Y", 999.9)
	bc.Mining()
	bc.Print()

	fmt.Printf("ME %.1f\n", bc.CalculateTotalAmount("my_blockchain_address"))
	fmt.Printf("C %.1f\n", bc.CalculateTotalAmount("C"))
	fmt.Printf("D %.1f\n", bc.CalculateTotalAmount("D"))
}
