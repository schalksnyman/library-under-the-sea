package logical

import (
	"context"
	repo "library-under-the-sea/services/library-repo/domain"
	library "library-under-the-sea/services/library/domain"
)

var _ library.Client = (*client)(nil)

func New(r repo.Client) *client {
	return &client{
		r: r,
	}
}

type client struct {
	r repo.Client
}

func (l *client) SaveBook(ctx context.Context, book library.Book) error {
	return nil
}

func (l *client) ListBooksByTitle(ctx context.Context, title string) ([]*library.Book, error) {
	return nil, nil
}

func (l *client) UpdateTitle(ctx context.Context, id int64, title string) error {
	return nil
}

func (l *client) DeleteBook(ctx context.Context, id int64) error {
	return nil
}

func (l *client) ListAll(ctx context.Context) ([]*library.Book, error) {
	return nil, nil
}
