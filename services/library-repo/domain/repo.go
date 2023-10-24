package domain

import (
	library "library-under-the-sea/services/library/domain"
)

type Client interface {
	Get(id string) (*library.Book, error)
	ListByTitle(title string) ([]*library.Book, error)
	ListAll() ([]*library.Book, error)
	Save(b library.Book) (string, error)
	Delete(id string) error
}
