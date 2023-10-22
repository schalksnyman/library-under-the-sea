package domain

import (
	library "library-under-the-sea/services/library/domain"
)

type Client interface {
	Get(ID int64) (*library.Book, error)
	Save(b library.Book) error
	Delete(ID int64) error
}
