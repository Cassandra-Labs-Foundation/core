package person

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Cassandra-Labs-Foundation/core/internal/service/person"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler provides HTTP handlers for person entity endpoints
type Handler struct {
	service person.Service
}

// NewHandler creates a new person handler
func NewHandler(service person.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Create handles the creation of a new person entity
// @Summary Create a new person entity
// @Description Create a new person entity with the provided information
// @Tags entities
// @Accept json
// @Produce json
// @Param input body person.CreatePersonInput true "Person creation input"
// @Success 201 {object} person.PersonOutput
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/entities/person [post]
func (h *Handler) Create(c *gin.Context) {
    var input person.CreatePersonInput
    if err := c.ShouldBindJSON(&input); err != nil {
        log.Printf("Error binding JSON: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    output, err := h.service.Create(c.Request.Context(), input)
    if err != nil {
        // Log the detailed error
        log.Printf("Error creating person entity: %v", err)
        
        if errors.Is(err, person.ErrInvalidPerson) {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person entity", "details": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, output)
}

// Get handles retrieving a person entity by ID
// @Summary Get a person entity by ID
// @Description Get a person entity by its unique identifier
// @Tags entities
// @Produce json
// @Param id path string true "Person ID"
// @Success 200 {object} person.PersonOutput
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/entities/person/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	output, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, person.ErrPersonNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Person entity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get person entity"})
		return
	}

	c.JSON(http.StatusOK, output)
}

// Update handles updating an existing person entity
// @Summary Update a person entity
// @Description Update an existing person entity with the provided information
// @Tags entities
// @Accept json
// @Produce json
// @Param id path string true "Person ID"
// @Param input body person.UpdatePersonInput true "Person update input"
// @Success 200 {object} person.PersonOutput
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/entities/person/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input person.UpdatePersonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.service.Update(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, person.ErrPersonNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Person entity not found"})
			return
		}
		if errors.Is(err, person.ErrInvalidPerson) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update person entity"})
		return
	}

	c.JSON(http.StatusOK, output)
}

// List handles retrieving a paginated list of person entities
// @Summary List person entities
// @Description Get a paginated list of person entities
// @Tags entities
// @Produce json
// @Param limit query int false "Limit (default 10, max 100)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {array} person.PersonOutput
// @Failure 500 {object} map[string]string
// @Router /api/v1/entities/person [get]
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list person entities"})
		return
	}

	c.JSON(http.StatusOK, outputs)
}