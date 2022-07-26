package domain

import (
	"html"
	"strings"
	"time"
)

type Comment struct {
	ID        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	PostID    uint64     `gorm:"size:100;not null;" json:"post_id"`
	Content   string     `gorm:"text;not null;" json:"content"`
	Creator   string     `gorm:"size:100;not null;" json:"creator"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (c *Comment) BeforeSave() {
	c.Content = html.EscapeString(strings.TrimSpace(c.Content))
}

func (c *Comment) Prepare() {
	c.Content = html.EscapeString(strings.TrimSpace(c.Content))
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Comment) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)

	switch strings.ToLower(action) {
	case "update":
		if c.Content == "" || c.Content == "null" {
			errorMessages["content_required"] = "content is required"
		}
	default:
		if c.Content == "" || c.Content == "null" {
			errorMessages["content_required"] = "content is required"
		}
		if c.Creator == "" || c.Creator == "null" {
			errorMessages["creator_required"] = "creator is required"
		}
	}
	return errorMessages
}
