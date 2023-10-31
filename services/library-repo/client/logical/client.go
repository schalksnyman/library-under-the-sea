package logical

import (
	"context"
	"library-under-the-sea/services/library-repo/client/ops"
	"library-under-the-sea/services/library-repo/client/ops/db"
	repo "library-under-the-sea/services/library-repo/domain"
	library "library-under-the-sea/services/library/domain"
)

var _ repo.Client = (*client)(nil)

func New(connectString string, dbName string) *client {
	dbHandler := db.NewMongoClient(connectString, dbName)

	return &client{
		d: dbHandler,
	}
}

type client struct {
	d repo.DBHandler
}

func (c *client) Get(ctx context.Context, id string) (*library.Book, error) {
	return ops.Get(ctx, id, c.d)
}

func (c *client) ListByTitle(ctx context.Context, title string) ([]*library.Book, error) {
	return ops.ListByTitle(ctx, title, c.d)
}

func (c *client) ListAll(ctx context.Context) ([]*library.Book, error) {
	return ops.ListAll(ctx, c.d)
}

func (c *client) Save(ctx context.Context, book library.Book) (string, error) {
	return ops.Save(ctx, book, c.d)
}

func (c *client) Delete(ctx context.Context, id string) error {
	return ops.Delete(ctx, id, c.d)
}
