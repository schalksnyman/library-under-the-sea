package grpc

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	library "library-under-the-sea/services/library/domain"
	"library-under-the-sea/services/library/librarypb"
	pb "library-under-the-sea/services/library/librarypb"
	"testing"
)

var addr = flag.String("library_address", "", "host:port of library gRPC service")

var _ library.Client = (*client)(nil)

type client struct {
	rpcConn   *grpc.ClientConn
	rpcClient pb.LibraryClient
}

func IsEnabled() bool {
	return *addr != ""
}

func New() (*client, error) {
	var c client
	var err error
	c.rpcConn, err = grpctls.NewClientAutoDiscoverReader(*addr, "authenticator", 7088, 30058)
	if err != nil {
		return nil, err
	}

	c.rpcClient = pb.NewLibraryClient(c.rpcConn)

	return &c, nil
}

func NewForTesting(_ *testing.T, conn *grpc.ClientConn) *client {
	return &client{
		rpcConn:   conn,
		rpcClient: pb.NewLibraryClient(conn),
	}
}

func (c *client) FindBook(ctx context.Context, id string) (*library.Book, error) {
	res, err := c.rpcClient.FindBook(ctx, &pb.FindBookRequest{Id: id})
	if err != nil {
		return nil, err
	}

	book := library.Book{
		Title:       res.Book.Title,
		Author:      res.Book.Author,
		Edition:     res.Book.Edition,
		PublishDate: res.Book.PublishDate.AsTime(),
	}

	return &book, nil
}

func (c *client) ListBooksByTitle(ctx context.Context, title string) ([]*library.Book, error) {
	res, err := c.rpcClient.ListBooksByTitle(ctx, &pb.ListBooksByTitleRequest{Title: title})
	if err != nil {
		return nil, err
	}

	var books []*library.Book
	for _, bookpb := range res.Books {
		book := &library.Book{
			Title:       bookpb.Title,
			Author:      bookpb.Author,
			Edition:     bookpb.Edition,
			PublishDate: bookpb.PublishDate.AsTime(),
		}

		books = append(books, book)
	}

	return books, nil
}

func (c *client) ListAll(ctx context.Context) ([]*library.Book, error) {
	res, err := c.rpcClient.ListAll(ctx, &pb.ListAllRequest{})
	if err != nil {
		return nil, err
	}

	var books []*library.Book
	for _, bookpb := range res.Books {
		book := &library.Book{
			Title:       bookpb.Title,
			Author:      bookpb.Author,
			Edition:     bookpb.Edition,
			PublishDate: bookpb.PublishDate.AsTime(),
		}

		books = append(books, book)
	}

	return books, nil
}

func (c *client) SaveBook(ctx context.Context, b library.Book) (string, error) {
	bookpb := librarypb.Book{
		Title:       b.Title,
		Author:      b.Author,
		Edition:     b.Edition,
		PublishDate: timestamppb.New(b.PublishDate),
	}

	res, err := c.rpcClient.SaveBook(ctx, &pb.SaveBookRequest{Book: &bookpb})
	if err != nil {
		return "", err
	}

	return res.Id, nil
}

func (c *client) UpdateTitle(ctx context.Context, id string, title string) error {
	_, err := c.rpcClient.UpdateTitle(ctx, &pb.UpdateTitleRequest{Id: id, Title: title})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) DeleteBook(ctx context.Context, id string) error {
	_, err := c.rpcClient.DeleteBook(ctx, &pb.DeleteBookRequest{Id: id})
	if err != nil {
		return err
	}

	return nil
}
