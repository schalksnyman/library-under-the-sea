package domain

import (
	"context"
	"time"
)

type Book struct {
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Edition     string    `json:"edition"`
	PublishDate time.Time `json:"publish_date"`
}

type Client interface {
	FindBook(ctx context.Context, id string) (*Book, error)
	ListBooksByTitle(ctx context.Context, title string) ([]*Book, error)
	ListAll(ctx context.Context) ([]*Book, error)
	SaveBook(ctx context.Context, book Book) (string, error)
	UpdateTitle(ctx context.Context, id string, title string) error
	DeleteBook(ctx context.Context, id string) error
}
