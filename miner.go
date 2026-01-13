package main

import (
	"fmt"
	"time"
)

// StartMiner starts a goroutine that mines a block every intervalSeconds.
// minerAddr is used as the header.Miner. difficulty is number of leading zero hex chars.
func StartMiner(bc *Blockchain, interval time.Duration, minerAddr string, difficulty uint64) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			// Collect pending txs
			txs := bc.PendingTransactions()

			parent := bc.LastBlock()
			header := BlockHeader{
				ParentHash: parent.Hash,
				Miner:      minerAddr,
				Number:     parent.Header.Number + 1,
				Time:       time.Now().Unix(),
				Difficulty: difficulty,
				Nonce:      0,
				StateRoot:  "",
			}

			fmt.Printf("Miner: tentando minerar bloco %d com %d txs...\n", header.Number, len(txs))

			hash, nonce, ok := validProof(header, txs, difficulty)
			if !ok {
				fmt.Printf("Miner: falha ao encontrar proof para bloco %d\n", header.Number)
				continue
			}
			header.Nonce = nonce
			b := &Block{Header: header, Transactions: txs, Hash: hash}
			bc.AddBlock(b)
			fmt.Printf("Miner: bloco %d minerado! hash=%s nonce=%d\n", header.Number, hash[:16], nonce)
		}
	}()
}
