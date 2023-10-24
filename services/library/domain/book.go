package domain

import (
	"time"
)

type Book struct {
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Edition     string    `json:"edition"`
	PublishDate time.Time `json:"publish_date"`
}

type Client interface {
	FindBook(id string) (*Book, error)
	SaveBook(book Book) (string, error)
	ListBooksByTitle(title string) ([]*Book, error)
	UpdateTitle(id string, title string) error
	DeleteBook(id string) error
	ListAll() ([]*Book, error)
}
