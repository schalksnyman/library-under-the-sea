package server

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"library-under-the-sea/services/library-repo/client/ops"
	"library-under-the-sea/services/library-repo/client/ops/db"
	libraryrepo "library-under-the-sea/services/library-repo/domain"
	pb "library-under-the-sea/services/library-repo/libraryrepopb"
	library "library-under-the-sea/services/library/domain"
	"library-under-the-sea/services/library/librarypb"
)

var _ pb.LibraryRepoServer = (*Server)(nil)

// Server implements the library grpc server.
type Server struct {
	d libraryrepo.DBHandler
	//writer pb.LibraryRepoClient
}

// New returns a new server instance.
func New(connectString string, dbName string) *Server {
	dbHandler := db.NewMongoClient(connectString, dbName)
	//var libraryRepoWriter pb.LibraryRepoClient
	//if writerConn != nil {
	//	libraryRepoWriter = pb.NewLibraryRepoClient(writerConn)
	//}

	return &Server{
		d: dbHandler,
		//writer: libraryRepoWriter,
	}
}

func (srv *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	book, err := ops.Get(ctx, req.Id, srv.d)
	if err != nil {
		return nil, err
	}

	bookpb := &librarypb.Book{
		Title:       book.Title,
		Author:      book.Author,
		Edition:     book.Edition,
		PublishDate: timestamppb.New(book.PublishDate),
	}

	return &pb.GetResponse{Book: bookpb}, err
}

func (srv *Server) ListByTitle(ctx context.Context, req *pb.ListByTitleRequest) (*pb.ListByTitleResponse, error) {
	books, err := ops.ListByTitle(ctx, req.Title, srv.d)
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

	return &pb.ListByTitleResponse{Books: bookspb}, err
}

func (srv *Server) ListAll(ctx context.Context, req *pb.ListAllRequest) (*pb.ListAllResponse, error) {
	books, err := ops.ListAll(ctx, srv.d)
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

func (srv *Server) Save(ctx context.Context, req *pb.SaveRequest) (*pb.SaveResponse, error) {
	book := library.Book{
		Title:       req.Book.Title,
		Author:      req.Book.Author,
		Edition:     req.Book.Edition,
		PublishDate: req.Book.PublishDate.AsTime(),
	}

	id, err := ops.Save(ctx, book, srv.d)
	if err != nil {
		return nil, err
	}

	return &pb.SaveResponse{Id: id}, err
}

func (srv *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := ops.Delete(ctx, req.Id, srv.d)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteResponse{}, err
}
