package ledger

import (
	"net/http"
	"strconv"

	"github.com/Cassandra-Labs-Foundation/core/internal/service/ledger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service ledger.Service
}

func NewHandler(service ledger.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateAccountHandler handles account creation.
func (h *Handler) CreateAccountHandler(c *gin.Context) {
	// Accept initialBalance as query parameter (default to 0)
	balanceStr := c.DefaultQuery("balance", "0")
	balance, err := strconv.ParseInt(balanceStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid balance parameter"})
		return
	}

	accountID, err := h.service.CreateAccount(c.Request.Context(), balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"account_id": accountID})
}

// TransferHandler handles fund transfers between accounts.
func (h *Handler) TransferHandler(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	amountStr := c.Query("amount")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil || amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount parameter"})
		return
	}

	err = h.service.TransferFunds(c.Request.Context(), from, to, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transfer funds", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}