package business

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Cassandra-Labs-Foundation/core/internal/service/business"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service business.Service
}

func NewHandler(service business.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var input business.CreateBusinessInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	output, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		log.Printf("Error creating business entity: %v", err)
		if errors.Is(err, business.ErrInvalidBusiness) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create business entity", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, output)
}

func (h *Handler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	output, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, business.ErrBusinessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Business entity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get business entity"})
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var input business.UpdateBusinessInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	output, err := h.service.Update(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, business.ErrBusinessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Business entity not found"})
			return
		}
		if errors.Is(err, business.ErrInvalidBusiness) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update business entity"})
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *Handler) List(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	outputs, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list business entities"})
		return
	}
	c.JSON(http.StatusOK, outputs)
}