package ops

import (
	repo "library-under-the-sea/services/library-repo/domain"
	"log"
)

// UpdateTitle finds the book by ID, updates the title and saves the book with the new title
func UpdateTitle(id int64, title string, r repo.Client) error {
	book, err := r.Get(id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	bookCopy := *book
	bookCopy.Title = title

	_, err = r.Save(bookCopy)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
