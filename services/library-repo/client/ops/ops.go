package ops

import (
	"context"
	repo "library-under-the-sea/services/library-repo/domain"
	library "library-under-the-sea/services/library/domain"
)

func Get(ctx context.Context, id string, d repo.DBHandler) (*library.Book, error) {
	return d.Get(ctx, id)
}

func ListByTitle(ctx context.Context, title string, d repo.DBHandler) ([]*library.Book, error) {
	return d.ListByTitle(ctx, title)
}

func ListAll(ctx context.Context, d repo.DBHandler) ([]*library.Book, error) {
	return d.ListAll(ctx)
}

func Save(ctx context.Context, book library.Book, d repo.DBHandler) (string, error) {
	return d.Save(ctx, book)
}

func Delete(ctx context.Context, id string, d repo.DBHandler) error {
	return d.Delete(ctx, id)
}
