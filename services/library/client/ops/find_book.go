package ops

import (
	"context"
	repo "library-under-the-sea/services/library-repo/domain"
	library "library-under-the-sea/services/library/domain"
	"log"
)

// FindBook returns a book at a given id
func FindBook(ctx context.Context, id string, r repo.Client) (*library.Book, error) {
	book, err := r.Get(ctx, id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return book, nil
}
