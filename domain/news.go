package domain

import "time"

type News struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	Slug      string          `json:"slug"`
	Status    string          `json:"status"`
	Content   string          `json:"content"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Topics    []NewsTopicList `json:"topics"`
}

type NewsUpdate struct {
	ID        string      `json:"id"`
	Title     string      `json:"title"`
	Slug      string      `json:"slug"`
	Status    string      `json:"status"`
	Content   string      `json:"content"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Topics    []NewsTopic `json:"topics"`
}

type NewsTopicList struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NewsTopic struct {
	ID        string    `json:"id"`
	TopicId   string    `json:"topic_id"`
	NewsId    string    `json:"news_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateNewsRequest struct {
	Title   string      `json:"title" validate:"required"`
	Slug    string      `json:"slug"`
	Status  string      `json:"status" validate:"required"`
	Content string      `json:"content" validate:"required"`
	Topic   []NewsTopic `json:"topics"`
	//Password string `json:"password" validate:"required,password"`
}
type CreateNewsTopicRequest struct {
	TopicId string `json:"topic_id" validate:"required"`
	NewsId  string `json:"news_id" validate:"required"`
	//Password string `json:"password" validate:"required,password"`
}
type NewsTopicNew struct {
	TopicId string `json:"topic_id"`
}
type UpdateNewsRequest struct {
	Title   string         `json:"title" validate:"required"`
	Slug    string         `json:"slug"`
	Status  string         `json:"status" validate:"required"`
	Content string         `json:"content" validate:"required"`
	Topic   []NewsTopicNew `json:"topics"`
}

type NewsFilter struct {
	Search string `json:"search" query:"search"`
}
