package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/priince9381/irm_backend/internal/models"
	"github.com/priince9381/irm_backend/internal/repository"
	"github.com/priince9381/irm_backend/internal/utils"
)

type AuthHandler struct {
	db        *repository.ElasticsearchDB
	secretKey string
}

func NewAuthHandler(db *repository.ElasticsearchDB, secretKey string) *AuthHandler {
	return &AuthHandler{
		db:        db,
		secretKey: secretKey,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	exists, err := h.db.Exists("users", map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"Email": req.Email,
			},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	// Create new user
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Save user to Elasticsearch
	userJSON, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user data"})
		return
	}

	_, err = h.db.Index("users", string(userJSON))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	result, err := h.db.Search("users", map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"Email": req.Email,
			},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	var user models.User
	if err := json.Unmarshal([]byte(result[0]), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user data"})
		return
	}

	// Check password
	if err := user.CheckPassword(req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email, h.secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
