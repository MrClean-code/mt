package service

import (
	"MTBlockchain/pkg/model"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

type BlockchainService struct {
	blockchain *model.Blockchain
}

func NewBlockchainService() *BlockchainService {
	bc := &model.Blockchain{
		Chain:               []model.Block{},
		CurrentTransactions: []model.Transaction{},
	}
	bc.Chain = append(bc.Chain, createGenesisBlock())
	return &BlockchainService{blockchain: bc}
}

func createGenesisBlock() model.Block {
	return model.Block{
		Index:        1,
		Timestamp:    time.Now().Unix(),
		Transactions: []model.Transaction{},
		Proof:        100,
		PreviousHash: "1",
	}
}

func (s *BlockchainService) Mine() model.Block {
	lastBlock := s.blockchain.Chain[len(s.blockchain.Chain)-1]
	lastProof := lastBlock.Proof
	proof := proofOfWork(lastProof)

	s.blockchain.CurrentTransactions = append(s.blockchain.CurrentTransactions, model.Transaction{
		ID:        uuid.New().String(),
		Sender:    "0",
		Recipient: "node_identifier",
		Amount:    1,
	})

	block := model.Block{
		Index:        len(s.blockchain.Chain) + 1,
		Timestamp:    time.Now().Unix(),
		Transactions: s.blockchain.CurrentTransactions,
		Proof:        proof,
		PreviousHash: hash(lastBlock),
	}

	s.blockchain.CurrentTransactions = []model.Transaction{}
	s.blockchain.Chain = append(s.blockchain.Chain, block)

	return block
}

func (s *BlockchainService) NewTransaction(sender, recipient string, amount int) int {
	transaction := model.Transaction{
		ID:        uuid.New().String(),
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
	s.blockchain.CurrentTransactions = append(s.blockchain.CurrentTransactions, transaction)
	return s.blockchain.Chain[len(s.blockchain.Chain)-1].Index + 1
}

func (s *BlockchainService) FullChain() *model.Blockchain {
	return s.blockchain
}

func (s *BlockchainService) GetTransactionByID(id string) *model.Transaction {
	for _, block := range s.blockchain.Chain {
		for _, transaction := range block.Transactions {
			if transaction.ID == id {
				return &transaction
			}
		}
	}
	for _, transaction := range s.blockchain.CurrentTransactions {
		if transaction.ID == id {
			return &transaction
		}
	}
	return nil
}

func proofOfWork(lastProof int) int {
	proof := 0
	for !validProof(lastProof, proof) {
		proof++
	}
	return proof
}

func validProof(lastProof, proof int) bool {
	guess := strconv.Itoa(lastProof) + strconv.Itoa(proof)
	guessHash := sha256.Sum256([]byte(guess))
	return strings.HasPrefix(hex.EncodeToString(guessHash[:]), "0000")
}

func hash(block model.Block) string {
	blockBytes, _ := json.Marshal(block)
	hash := sha256.Sum256(blockBytes)
	return hex.EncodeToString(hash[:])
}
