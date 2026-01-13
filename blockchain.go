package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// BlockHeader inspirado na estrutura do Ethereum (simplificado)
type BlockHeader struct {
	ParentHash string `json:"parent_hash"`
	Miner      string `json:"miner"`
	Number     uint64 `json:"number"`
	Time       int64  `json:"time"`
	Difficulty uint64 `json:"difficulty"`
	Nonce      uint64 `json:"nonce"`
	Proof      []byte `json:"proof"`
	StateRoot  string `json:"state_root"` // placeholder
}

type Transaction struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint64 `json:"value"`
	Data  string `json:"data,omitempty"`
}

type Block struct {
	Header       BlockHeader   `json:"header"`
	Transactions []Transaction `json:"transactions"`
	Hash         string        `json:"hash"`
}

type Blockchain struct {
	mu         sync.Mutex
	blocks     []*Block
	pendingTxs []Transaction
	path       string
}

func NewBlockchain() *Blockchain {
	// ensure tmp directories exist
	ensureTmp()

	bc := &Blockchain{path: "./tmp/blocks"}
	// create genesis
	genesis := &Block{
		Header: BlockHeader{
			ParentHash: "",
			Miner:      "genesis",
			Number:     0,
			Time:       time.Now().Unix(),
			Difficulty: 1,
			Nonce:      0,
			Proof:      nil,
			StateRoot:  "",
		},
		Transactions: nil,
	}
	genesis.Hash = calcBlockHash(genesis)
	bc.blocks = []*Block{genesis}
	// persist genesis
	_ = saveBlockToFile(bc.path, genesis)
	return bc
}

func (bc *Blockchain) AddTransaction(tx Transaction) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.pendingTxs = append(bc.pendingTxs, tx)
}

func (bc *Blockchain) PendingTransactions() []Transaction {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	// return a copy
	copyTxs := make([]Transaction, len(bc.pendingTxs))
	copy(copyTxs, bc.pendingTxs)
	return copyTxs
}

func (bc *Blockchain) LastBlock() *Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	return bc.blocks[len(bc.blocks)-1]
}

func (bc *Blockchain) AddBlock(b *Block) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	// remove included txs from pending (simple approach: clear all)
	bc.pendingTxs = nil
	bc.blocks = append(bc.blocks, b)
	// persist block
	_ = saveBlockToFile(bc.path, b)
}

func (bc *Blockchain) Chain() []*Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	copyBlocks := make([]*Block, len(bc.blocks))
	copy(copyBlocks, bc.blocks)
	return copyBlocks
}

func calcBlockHash(b *Block) string {
	// marshal header + txs
	data, _ := json.Marshal(struct {
		Header       BlockHeader   `json:"header"`
		Transactions []Transaction `json:"transactions"`
	}{b.Header, b.Transactions})
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// Simple PoW: hash(header+nonce) interpreted as big-endian hex must have a number of leading zeros
func validProof(header BlockHeader, txs []Transaction, difficulty uint64) (string, uint64, bool) {
	// Try nonces until a hash with required zeros is found
	var nonce uint64
	for nonce = 0; nonce < 1<<20; nonce++ { // limit
		header.Nonce = nonce
		payload, _ := json.Marshal(struct {
			Header       BlockHeader   `json:"header"`
			Transactions []Transaction `json:"transactions"`
		}{header, txs})
		h := sha256.Sum256(payload)
		hexh := hex.EncodeToString(h[:])
		// check leading zeros (each zero nibble = 4 bits)
		need := int(difficulty)
		ok := true
		for i := 0; i < need; i++ {
			if i >= len(hexh) || hexh[i] != '0' {
				ok = false
				break
			}
		}
		if ok {
			return hexh, nonce, true
		}
	}
	return "", 0, false
}

// simple file storage: save block as JSON to tmp/blocks/<number>.json
func ensureTmp() error {
	_ = os.MkdirAll("./tmp/blocks", 0755)
	_ = os.MkdirAll("./tmp/wallets", 0755)
	return nil
}

func saveBlockToFile(path string, b *Block) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%s/block-%d.json", path, b.Header.Number)
	return os.WriteFile(filename, data, 0644)
}
