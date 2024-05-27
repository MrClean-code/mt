package handler

import (
	"MTBlockchain/pkg/model"
	"MTBlockchain/pkg/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

var blockchainService *service.BlockchainService

func init() {
	blockchainService = service.NewBlockchainService()
}

func RegisterRoutes(router *gin.Engine) {
	router.GET("/mine", mine)
	router.POST("/transactions/new", newTransaction)
	router.GET("/chain", fullChain)
	router.GET("/transactions", getTransaction)
}

func mine(c *gin.Context) {
	block := blockchainService.Mine()
	c.JSON(http.StatusOK, gin.H{
		"message":       "New Block Forged",
		"index":         block.Index,
		"transactions":  block.Transactions,
		"proof":         block.Proof,
		"previous_hash": block.PreviousHash,
	})
}

func newTransaction(c *gin.Context) {
	var t model.Transaction
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	index := blockchainService.NewTransaction(t.Sender, t.Recipient, t.Amount)
	c.JSON(http.StatusCreated, gin.H{"message": "Transaction will be added to Block " + strconv.Itoa(index)})
}

func fullChain(c *gin.Context) {
	chain := blockchainService.FullChain()
	c.JSON(http.StatusOK, gin.H{
		"chain":  chain.Chain,
		"length": len(chain.Chain),
	})
}

func getTransaction(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		log.Fatal("ID is empty") // Прерываем выполнение программы и выводим сообщение об ошибке
	}

	log.Println("id =", id)
	transaction := blockchainService.GetTransactionByID(id)
	if transaction == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	c.JSON(http.StatusOK, transaction)
}
