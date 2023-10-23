package ops

import (
	repo "library-under-the-sea/services/library-repo/domain"
	library "library-under-the-sea/services/library/domain"
	"log"
)

func ListAll(r repo.Client) ([]*library.Book, error) {
	books, err := r.ListAll()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return books, nil
}
