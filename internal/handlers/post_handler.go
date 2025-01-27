package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/priince9381/irm_backend/internal/models"
	"gorm.io/gorm"
)

type PostHandler struct {
	db        *gorm.DB
	uploadDir string
}

type CreatePostRequest struct {
	Content      string   `json:"content" binding:"required"`
	MediaURLs    []string `json:"media_urls"`
	Platforms    []string `json:"platforms" binding:"required"`
	ScheduledFor string   `json:"scheduled_for"`
}

func NewPostHandler(db *gorm.DB, uploadDir string) *PostHandler {
	return &PostHandler{
		db:        db,
		uploadDir: uploadDir,
	}
}

func (h *Handler) CreatePost(c *gin.Context) {
	// Parse multipart form with 10 MB max memory
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	// Get form values
	content := c.PostForm("content")
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		return
	}

	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	platforms := c.PostFormArray("platforms")
	if len(platforms) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one platform is required"})
		return
	}

	links := c.PostFormArray("links")

	// Parse scheduled time if provided
	var scheduledTime time.Time
	if timeStr := c.PostForm("scheduled_time"); timeStr != "" {
		var err error
		scheduledTime, err = time.Parse(time.RFC3339, timeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduled time format"})
			return
		}
	}

	status := c.PostForm("status")
	if status == "" {
		status = "draft"
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Handle file uploads
	form, _ := c.MultipartForm()
	files := form.File["files"]

	var mediaFiles []models.Media
	if len(files) > 0 {
		uploadDir := "uploads/"
		// Create uploads directory if it doesn't exist
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}

		for _, file := range files {
			// Generate unique filename
			filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
			filepath := uploadDir + filename

			// Save the file
			if err := c.SaveUploadedFile(file, filepath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
				return
			}

			// Determine media type
			mediaType := "image"
			if strings.HasPrefix(file.Header.Get("Content-Type"), "video/") {
				mediaType = "video"
			}

			mediaFiles = append(mediaFiles, models.Media{
				URL:      filepath,
				Type:     mediaType,
				FileName: filename,
			})
		}
	}

	post := &models.Post{
		Title:         title,
		Content:       content,
		UserID:        userID,
		Platforms:     platforms,
		MediaFiles:    mediaFiles,
		Links:         links,
		ScheduledTime: scheduledTime,
		Status:        status,
	}

	if err := h.db.CreatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"post":    post,
	})
}

func (h *Handler) GetPosts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	posts, err := h.db.GetUserPosts(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func (h *Handler) UpdatePost(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement post update logic
	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

func (h *Handler) DeletePost(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}

	if err := h.db.DeletePost(postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
