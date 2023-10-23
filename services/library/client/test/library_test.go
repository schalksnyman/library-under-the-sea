package test

import (
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
	repo_mock "library-under-the-sea/services/library-repo/client/mocks"
	"library-under-the-sea/services/library/client/logical"
	library "library-under-the-sea/services/library/domain"
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
