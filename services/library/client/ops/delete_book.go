package ops

import (
	repo "library-under-the-sea/services/library-repo/domain"
	"log"
)

func DeleteBook(id int64, r repo.Client) error {
	err := r.Delete(id)
	if err != nil {
		log.Println(err.Error())
	}

	return nil
}
