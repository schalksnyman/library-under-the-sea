package domain

import (
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
	FindBook(id int64) (*Book, error)
	SaveBook(book Book) (int64, error)
	ListBooksByTitle(title string) ([]*Book, error)
	UpdateTitle(id int64, title string) error
	DeleteBook(id int64) error
	ListAll() ([]*Book, error)
}
