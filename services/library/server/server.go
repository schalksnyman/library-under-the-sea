package server

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	libraryrepo "library-under-the-sea/services/library-repo/domain"
	"library-under-the-sea/services/library/client/ops"
	library "library-under-the-sea/services/library/domain"
	"library-under-the-sea/services/library/librarypb"
	pb "library-under-the-sea/services/library/librarypb"
)

var _ pb.LibraryServer = (*Server)(nil)

// Server implements the library grpc server.
type Server struct {
	repo   libraryrepo.Client
	writer pb.LibraryClient
}

// New returns a new server instance.
func New(r libraryrepo.Client, writerConn *grpc.ClientConn) *Server {
	var libraryRepoWriter pb.LibraryClient
	if writerConn != nil {
		libraryRepoWriter = pb.NewLibraryClient(writerConn)
	}

	return &Server{
		repo:   r,
		writer: libraryRepoWriter,
	}
}

func (srv *Server) FindBook(ctx context.Context, req *pb.FindBookRequest) (*pb.FindBookResponse, error) {
	book, err := ops.FindBook(ctx, req.Id, srv.repo)
	if err != nil {
		return nil, err
	}

	bookpb := &librarypb.Book{
		Title:       book.Title,
		Author:      book.Author,
		Edition:     book.Edition,
		PublishDate: timestamppb.New(book.PublishDate),
	}

	return &pb.FindBookResponse{Book: bookpb}, err
}

func (srv *Server) ListBooksByTitle(ctx context.Context, req *pb.ListBooksByTitleRequest) (*pb.ListBooksByTitleResponse, error) {
	books, err := ops.ListBooksByTitle(ctx, req.Title, srv.repo)
	if err != nil {
		return nil, err
	}

	var bookspb []*librarypb.Book
	for _, book := range books {
		bookpb := &librarypb.Book{
			Title:       book.Title,
			Author:      book.Author,
			Edition:     book.Edition,
			PublishDate: timestamppb.New(book.PublishDate),
		}

		bookspb = append(bookspb, bookpb)
	}

	return &pb.ListBooksByTitleResponse{Books: bookspb}, err
}

func (srv *Server) ListAll(ctx context.Context, req *pb.ListAllRequest) (*pb.ListAllResponse, error) {
	books, err := ops.ListAll(ctx, srv.repo)
	if err != nil {
		return nil, err
	}

	var bookspb []*librarypb.Book
	for _, book := range books {
		bookpb := &librarypb.Book{
			Title:       book.Title,
			Author:      book.Author,
			Edition:     book.Edition,
			PublishDate: timestamppb.New(book.PublishDate),
		}

		bookspb = append(bookspb, bookpb)
	}

	return &pb.ListAllResponse{Books: bookspb}, err
}

func (srv *Server) SaveBook(ctx context.Context, req *pb.SaveBookRequest) (*pb.SaveBookResponse, error) {
	book := library.Book{
		Title:       req.Book.Title,
		Author:      req.Book.Author,
		Edition:     req.Book.Edition,
		PublishDate: req.Book.PublishDate.AsTime(),
	}

	id, err := ops.SaveBook(ctx, book, srv.repo)
	if err != nil {
		return nil, err
	}

	return &pb.SaveBookResponse{Id: id}, err
}

func (srv *Server) UpdateTitle(ctx context.Context, req *pb.UpdateTitleRequest) (*pb.UpdateTitleResponse, error) {
	err := ops.UpdateTitle(ctx, req.Id, req.Title, srv.repo)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateTitleResponse{}, err
}

func (srv *Server) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	err := ops.DeleteBook(ctx, req.Id, srv.repo)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteBookResponse{}, err
}
