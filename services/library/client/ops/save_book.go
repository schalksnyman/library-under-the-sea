package ops

import (
	repo "library-under-the-sea/services/library-repo/domain"
	library "library-under-the-sea/services/library/domain"
	"log"
)

// SaveBook saves the book to the library
func SaveBook(book library.Book, r repo.Client) (int64, error) {
	id, err := r.Save(book)
	if err != nil {
		log.Println(err.Error())
		return -1, err
	}

	return id, nil
}
