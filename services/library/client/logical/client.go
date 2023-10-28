package logical

import (
	"context"
	repo "library-under-the-sea/services/library-repo/domain"
	"library-under-the-sea/services/library/client/ops"
	library "library-under-the-sea/services/library/domain"
	"log"
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

func (c *client) FindBook(ctx context.Context, id string) (*library.Book, error) {
	book, err := ops.FindBook(ctx, id, c.r)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return book, nil
}

func (c *client) ListBooksByTitle(ctx context.Context, title string) ([]*library.Book, error) {
	books, err := ops.ListBooksByTitle(ctx, title, c.r)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return books, nil
}

func (c *client) ListAll(ctx context.Context) ([]*library.Book, error) {
	books, err := ops.ListAll(ctx, c.r)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return books, nil
}

func (c *client) SaveBook(ctx context.Context, book library.Book) (string, error) {
	id, err := ops.SaveBook(ctx, book, c.r)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return id, nil
}

func (c *client) UpdateTitle(ctx context.Context, id string, title string) error {
	err := ops.UpdateTitle(ctx, id, title, c.r)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (c *client) DeleteBook(ctx context.Context, id string) error {
	err := ops.DeleteBook(ctx, id, c.r)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
