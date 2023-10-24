package logical

import (
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

func (c *client) FindBook(id string) (*library.Book, error) {
	book, err := ops.FindBook(id, c.r)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return book, nil
}

func (c *client) SaveBook(book library.Book) (string, error) {
	id, err := ops.SaveBook(book, c.r)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return id, nil
}

func (c *client) ListBooksByTitle(title string) ([]*library.Book, error) {
	books, err := ops.ListBooksByTitle(title, c.r)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return books, nil
}

func (c *client) UpdateTitle(id string, title string) error {
	err := ops.UpdateTitle(id, title, c.r)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (c *client) DeleteBook(id string) error {
	err := ops.DeleteBook(id, c.r)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (c *client) ListAll() ([]*library.Book, error) {
	books, err := ops.ListAll(c.r)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return books, nil
}
