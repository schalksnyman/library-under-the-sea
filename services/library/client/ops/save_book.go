package ops

import (
	"context"
	repo "library-under-the-sea/services/library-repo/domain"
	library "library-under-the-sea/services/library/domain"
	"log"
)

// SaveBook saves the book to the library and returns the id
func SaveBook(ctx context.Context, book library.Book, r repo.Client) (string, error) {
	id, err := r.Save(book)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return id, nil
}
