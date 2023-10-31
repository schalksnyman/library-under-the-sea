package ops

import (
	"context"
	repo "library-under-the-sea/services/library-repo/domain"
	"log"
)

// UpdateTitle finds the book by ID, updates the title and saves the book with the new title
func UpdateTitle(ctx context.Context, id string, title string, r repo.Client) error {
	book, err := r.Get(ctx, id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	bookCopy := *book
	bookCopy.Title = title

	//TODO(ssnyman): Rather create update repo method
	_, err = r.Save(ctx, bookCopy)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
