package domain

import (
	library "library-under-the-sea/services/library/domain"
)

type Client interface {
	Get(id int64) (*library.Book, error)
	ListByTitle(title string) ([]*library.Book, error)
	ListAll() ([]*library.Book, error)
	Save(b library.Book) (int64, error)
	Delete(id int64) error
}
