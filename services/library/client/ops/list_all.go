package ops

import (
	"context"
	repo "library-under-the-sea/services/library-repo/domain"
	library "library-under-the-sea/services/library/domain"
	"log"
)

// ListAll returns all the books in the library
func ListAll(ctx context.Context, r repo.Client) ([]*library.Book, error) {
	books, err := r.ListAll()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return books, nil
}
