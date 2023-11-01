package agent

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"library-under-the-sea/services/library-api/handler"
	libraryrepoClientDev "library-under-the-sea/services/library-repo/client/dev"
	libraryClientDev "library-under-the-sea/services/library/client/dev"
	"log"
	"net"
	"net/http"
	"sync"
)

type Config struct {
	// HTTPAddr is the healthcheck address
	HTTPAddr string
	// Library Repo GRPC address
	LibraryRepoAddr string
	// Library GRPC address
	LibraryAddr string
	// Database connect string
	DBConnectString string
	// Database name
	DBName string
}

type Agent struct {
	Config Config

	server *Server

	shutdown     bool
	shutdowns    chan struct{}
	shutdownLock sync.Mutex
}

func New(config Config) (*Agent, error) {
	a := &Agent{
		Config:    config,
		shutdowns: make(chan struct{}),
	}
	setup := []func() error{
		a.setupHTTPServer,
	}
	for _, fn := range setup {
		if err := fn(); err != nil {
			return nil, err
		}
	}
	go a.serve()
	return a, nil
}

func (a *Agent) setupHTTPServer() error {
	log.Printf("Set up HTTP server\n")

	// Validate address
	addr, err := net.ResolveTCPAddr("tcp", a.Config.HTTPAddr)
	if err != nil {
		return err
	}

	if addr == nil {
		return errors.New("no address provided")
	}

	log.Printf("Create library repo client connection to GRPC server %s\n", a.Config.LibraryRepoAddr)
	libraryRepoClient, err := libraryrepoClientDev.New(a.Config.LibraryRepoAddr, a.Config.DBConnectString, a.Config.DBName)
	if err != nil {
		return errors.New(err.Error())
	}

	log.Printf("Create library client connection to GRPC server %s\n", a.Config.LibraryAddr)
	libraryClient, err := libraryClientDev.New(a.Config.LibraryAddr, libraryRepoClient)
	if err != nil {
		return errors.New(err.Error())
	}

	libraryAPIHandler := handler.New(libraryClient)

	r := mux.NewRouter()
	r.HandleFunc("/health", makeHealthCheckHandler()).Methods(http.MethodGet)
	r.HandleFunc("/library/add", libraryAPIHandler.Add).Methods(http.MethodPost)
	r.HandleFunc("/library/findall", libraryAPIHandler.FindAll).Methods(http.MethodPost)
	r.HandleFunc("/library/findbytitle", libraryAPIHandler.FindByTitle).Methods(http.MethodPost)
	r.HandleFunc("/library/findbook", libraryAPIHandler.FindBook).Methods(http.MethodPost)
	r.HandleFunc("/library/updatebooktitle", libraryAPIHandler.UpdateBookTitle).Methods(http.MethodPost)
	r.HandleFunc("/library/deletebook", libraryAPIHandler.DeleteBook).Methods(http.MethodPost)

	httpServer := &http.Server{
		Addr:    addr.String(),
		Handler: r,
	}
	httpServer.Addr = a.Config.HTTPAddr

	var srv Server
	srv.httpServer = httpServer

	a.server = &srv

	return nil
}

func makeHealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
	}
}

func (a *Agent) serve() error {
	if err := a.server.ServeForever(); err != nil {
		_ = a.Shutdown()
		return err
	}

	return nil
}

func (a *Agent) Shutdown() error {
	a.shutdownLock.Lock()
	defer a.shutdownLock.Unlock()
	if a.shutdown {
		return nil
	}
	a.shutdown = true
	close(a.shutdowns)

	shutdown := []func() error{
		func() error {
			a.server.Stop()
			return nil
		},
	}
	for _, fn := range shutdown {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

// Server wraps a gRPC server.
type Server struct {
	listener   net.Listener
	grpcServer *grpc.Server
	httpServer *http.Server
}

// Listener returns the server's net.Listener.
func (srv *Server) Listener() net.Listener {
	return srv.listener
}

// Stop stops the gRPC server.
func (srv *Server) Stop() {
	srv.httpServer.Shutdown(context.TODO())
}

// ServeForever listens for gRPC requests.
func (srv *Server) ServeForever() error {
	// Serve HTTP server
	log.Printf("Starting up HTTP server\n")
	errServerCh := make(chan error)
	go func() {
		err := srv.httpServer.ListenAndServe()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				errServerCh <- err
			}
		}
		errServerCh <- nil
	}()

	return <-errServerCh
}
