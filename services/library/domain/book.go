package domain

import (
	"context"
	"time"
)

type Book struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Edition     string    `json:"edition"`
	PublishDate time.Time `json:"publish_date"`
}

type Client interface {
	SaveBook(ctx context.Context, book Book) error
	ListBooksByTitle(ctx context.Context, title string) ([]*Book, error)
	UpdateTitle(ctx context.Context, id int64, title string) error
	DeleteBook(ctx context.Context, id int64) error
	ListAll(ctx context.Context) ([]*Book, error)
}
