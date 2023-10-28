package domain

import (
	"context"
	library "library-under-the-sea/services/library/domain"
)

type DBHandler interface {
	Get(ctx context.Context, id string) (*library.Book, error)
	ListByTitle(ctx context.Context, title string) ([]*library.Book, error)
	ListAll(ctx context.Context) ([]*library.Book, error)
	Save(ctx context.Context, b library.Book) (string, error)
	Delete(ctx context.Context, id string) error
}
