package models

import "time"

type (
	Post struct {
		ID        int       `json:"id"`
		Title     string    `json:"title"`
		Content   string    `json:"content" `
		Slug      string    `json:"slug" `
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	PostWithTag struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Tags    []Tag  `gorm:"many2many:post_tags;foreignKey:ID;joinForeignKey:IdPost;joinReferences:IdTag"`
	}
)

func (PostWithTag) TableName() string {
	return "posts"
}
