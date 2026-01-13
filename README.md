Block-go
========

Exemplo mínimo de blockchain inspirado na estrutura de bloco do Ethereum (muito simplificada).

Como rodar:

```bash
go run main.go blockchain.go miner.go api.go
```

Endpoints:

- POST /tx  - body JSON: {"from":"a","to":"b","value":100}
- GET /chain - retorna a cadeia em JSON

O minerador roda a cada 15 segundos por padrão e tenta encontrar um hash com N zeros no começo (difficulty).
