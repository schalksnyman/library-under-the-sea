package ops

import (
	repo "library-under-the-sea/services/library-repo/domain"
	library "library-under-the-sea/services/library/domain"
	"log"
)

// ListBooksByTitle returns a list of books filtered by title
func ListBooksByTitle(title string, r repo.Client) ([]*library.Book, error) {
	books, err := r.ListByTitle(title)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return books, nil
}
