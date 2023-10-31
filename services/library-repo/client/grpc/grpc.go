package grpc

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	repo "library-under-the-sea/services/library-repo/domain"
	pb "library-under-the-sea/services/library-repo/libraryrepopb"
	library "library-under-the-sea/services/library/domain"
	"library-under-the-sea/services/library/librarypb"
	"testing"
)

var addr = flag.String("library_repo_address", "", "host:port of library_repo gRPC service")

var _ repo.Client = (*client)(nil)

type client struct {
	rpcConn   *grpc.ClientConn
	rpcClient pb.LibraryRepoClient
}

func IsEnabled() bool {
	return *addr != ""
}

func New(connectString string, dbName string) (*client, error) {
	var c client
	var err error
	c.rpcConn, err = newGRPCConnection(*addr)
	if err != nil {
		return nil, err
	}

	c.rpcClient = pb.NewLibraryRepoClient(c.rpcConn)

	return &c, nil
}

func NewForTesting(_ *testing.T, conn *grpc.ClientConn) *client {
	return &client{
		rpcConn:   conn,
		rpcClient: pb.NewLibraryRepoClient(conn),
	}
}

func (c *client) Get(ctx context.Context, id string) (*library.Book, error) {
	res, err := c.rpcClient.Get(ctx, &pb.GetRequest{Id: id})
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

func (c *client) ListByTitle(ctx context.Context, title string) ([]*library.Book, error) {
	res, err := c.rpcClient.ListByTitle(ctx, &pb.ListByTitleRequest{Title: title})
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

func (c *client) Save(ctx context.Context, b library.Book) (string, error) {
	bookpb := librarypb.Book{
		Title:       b.Title,
		Author:      b.Author,
		Edition:     b.Edition,
		PublishDate: timestamppb.New(b.PublishDate),
	}

	res, err := c.rpcClient.Save(ctx, &pb.SaveRequest{Book: &bookpb})
	if err != nil {
		return "", err
	}

	return res.Id, nil
}

func (c *client) Delete(ctx context.Context, id string) error {
	_, err := c.rpcClient.Delete(ctx, &pb.DeleteRequest{Id: id})
	if err != nil {
		return err
	}

	return nil
}

func newGRPCConnection(addr string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption // No options yet

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
