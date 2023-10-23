package test

import (
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
	repo_mock "library-under-the-sea/services/library-repo/client/mocks"
	"library-under-the-sea/services/library/client/logical"
	library "library-under-the-sea/services/library/domain"
	"math/rand"
	"testing"
)

func TestLogical(t *testing.T) {
	suite.Run(t, new(TestLogicalSuite))
}

type TestLogicalSuite struct {
	clientSuite
}

func (s *TestLogicalSuite) SetupSuite() {
	s.T().Cleanup(func() {
		goleak.VerifyNone(s.T())
	})
}

type clientSuite struct {
	suite.Suite

	client library.Client
}

func (s *clientSuite) TestFindBook() {
	repoClient := new(repo_mock.Client).TSetup(s.T())

	f := fuzz.New()

	var bookExample library.Book
	f.Fuzz(&bookExample)

	repoClient.EXPECT().Get(mock.Anything).Return(&bookExample, nil)

	s.client = logical.New(repoClient)

	s.Run("Successfully Find Book", func() {
		book, err := s.client.FindBook(bookExample.Id)
		s.Require().NoError(err)
		s.Require().NotNil(book)
	})
}

func (s *clientSuite) TestSaveBook() {
	repoClient := new(repo_mock.Client).TSetup(s.T())
	repoClient.EXPECT().Save(mock.Anything).Return(1, nil)

	s.client = logical.New(repoClient)

	f := fuzz.New()

	var bookExample library.Book
	f.Fuzz(&bookExample)

	s.Run("Successfully Save Book", func() {
		id, err := s.client.SaveBook(bookExample)
		s.Require().NoError(err)
		s.Require().GreaterOrEqual(id, int64(0))
	})
}

func (s *clientSuite) TestListBooksByTitle() {
	repoClient := new(repo_mock.Client).TSetup(s.T())

	unicodeRanges := fuzz.UnicodeRanges{
		{First: 'a', Last: 'z'},
		{First: '0', Last: '9'},
	}
	f := fuzz.New().Funcs(unicodeRanges.CustomStringFuzzFunc())

	var bookExample1 library.Book
	var bookExample2 library.Book

	f.Fuzz(&bookExample1)
	f.Fuzz(&bookExample2)

	searchTitle := "Just for illustration. Not used."
	bookExample1.Title = searchTitle
	bookExample2.Title = searchTitle

	bookExamples := []*library.Book{
		&bookExample1,
		&bookExample2,
	}

	repoClient.EXPECT().ListByTitle(mock.Anything).Return(bookExamples, nil)

	s.client = logical.New(repoClient)

	s.Run("Successfully List Books for Title", func() {
		bookExamples, err := s.client.ListBooksByTitle(searchTitle)
		s.Require().NoError(err)
		s.Require().Len(bookExamples, 2)
	})
}

func (s *clientSuite) TestUpdateTitle() {
	repoClient := new(repo_mock.Client).TSetup(s.T())

	f := fuzz.New()

	var bookExample library.Book
	f.Fuzz(&bookExample)

	repoClient.EXPECT().Get(mock.Anything).Once().Return(&bookExample, nil)

	bookExampleCopy := bookExample
	newTitle := "Updated title"
	bookExampleCopy.Title = newTitle
	repoClient.EXPECT().Save(mock.Anything).Once().Return(bookExampleCopy.Id, nil)

	repoClient.EXPECT().Get(bookExampleCopy.Id).Once().Return(&bookExampleCopy, nil)

	s.client = logical.New(repoClient)

	s.Run("Successfully Update Title", func() {
		err := s.client.UpdateTitle(bookExample.Id, newTitle)
		s.Require().NoError(err)

		book, err := s.client.FindBook(bookExample.Id)
		s.Require().NoError(err)
		s.Require().Equal(newTitle, book.Title)
	})
}

func (s *clientSuite) TestDeleteBook() {
	repoClient := new(repo_mock.Client).TSetup(s.T())
	repoClient.EXPECT().Delete(mock.Anything).Return(nil)

	s.client = logical.New(repoClient)

	id := rand.Int63()

	s.Run("Successfully Save Book", func() {
		err := s.client.DeleteBook(id)
		s.Require().NoError(err)
	})
}

func (s *clientSuite) TestListAll() {
	repoClient := new(repo_mock.Client).TSetup(s.T())

	unicodeRanges := fuzz.UnicodeRanges{
		{First: 'a', Last: 'z'},
		{First: '0', Last: '9'},
	}
	f := fuzz.New().Funcs(unicodeRanges.CustomStringFuzzFunc())

	var bookExamples []*library.Book

	for i := 0; i < 10; i++ {
		var b library.Book
		f.Fuzz(&b)
		bookExamples = append(bookExamples, &b)
	}

	repoClient.EXPECT().ListAll().Return(bookExamples, nil)

	s.client = logical.New(repoClient)

	s.Run("Successfully List All Books", func() {
		bookExamples, err := s.client.ListAll()
		s.Require().NoError(err)
		s.Require().Len(bookExamples, 10)
	})
}
