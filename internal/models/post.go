package models

import (
	"time"
)

type Post struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Platforms     []string  `json:"platforms"`
	MediaFiles    []Media   `json:"media_files,omitempty"`
	Links         []string  `json:"links,omitempty"`
	ScheduledTime time.Time `json:"scheduled_time,omitempty"`
	Status        string    `json:"status"` // draft, scheduled, published
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserID        string    `json:"user_id"`
}

type Media struct {
	URL      string `json:"url"`
	Type     string `json:"type"` // image, video
	FileName string `json:"file_name"`
}

type CreatePostRequest struct {
	Title         string    `json:"title" binding:"required"`
	Content       string    `json:"content" binding:"required"`
	Platforms     []string  `json:"platforms" binding:"required"`
	Links         []string  `json:"links"`
	ScheduledTime time.Time `json:"scheduled_time"`
	Status        string    `json:"status" binding:"required,oneof=draft scheduled published"`
}

type UpdatePostRequest struct {
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Platforms     []string  `json:"platforms"`
	Links         []string  `json:"links"`
	ScheduledTime time.Time `json:"scheduled_time"`
	Status        string    `json:"status" binding:"oneof=draft scheduled published"`
}
