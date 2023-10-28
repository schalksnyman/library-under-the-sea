package ops

import (
	"context"
	repo "library-under-the-sea/services/library-repo/domain"
	"log"
)

// DeleteBook deletes a book at a given id
func DeleteBook(ctx context.Context, id string, r repo.Client) error {
	err := r.Delete(id)
	if err != nil {
		log.Println(err.Error())
	}

	return nil
}
