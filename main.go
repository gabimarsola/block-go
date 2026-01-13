package main

import (
	chain "block-go/structure/blockchain"
	"block-go/structure/miner"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Inicializa blockchain
	chain := chain.InitBlockchain()
	defer chain.Database.Close()

	miner.Mine()

	// Ensure tmp and wallets
	//_, _ = walletpkg.CreateWallets()

	// Inicia minerador que roda a cada 15s (executor simples de PoW)
	//StartMiner(bc, 15*time.Second, "miner1", 3)

	// Inicia API HTTP
	mux := http.NewServeMux()
	registerHandlers(mux, chain)

	addr := ":8080"
	fmt.Printf("Servidor HTTP rodando em %s - minerador a cada 15s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
