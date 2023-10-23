package ops

import (
	repo "library-under-the-sea/services/library-repo/domain"
	library "library-under-the-sea/services/library/domain"
	"log"
)

// FindBook returns a book at a given id
func FindBook(id int64, r repo.Client) (*library.Book, error) {
	book, err := r.Get(id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return book, nil
}
