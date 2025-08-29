package domain

import (
	"time"
)

type Topic struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTopicRequest struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug"`
	//Password string `json:"password" validate:"required,password"`
}

type UpdateTopicRequest struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug"`
}

type TopicFilter struct {
	Search string `json:"search" query:"search"`
}
