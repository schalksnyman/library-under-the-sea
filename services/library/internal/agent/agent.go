package agent

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	libraryrepoClientDev "library-under-the-sea/services/library-repo/client/dev"
	"library-under-the-sea/services/library/librarypb"
	library_server "library-under-the-sea/services/library/server"
)

type Config struct {
	// HTTPAddr is the healthcheck address
	HTTPAddr string
	// GRPCAddr is the GRPC address
	GRPCAddr string
	// Library repo GRPC address
	LibraryRepoAddr string
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
		a.setupGRPCServer,
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

func (a *Agent) setupGRPCServer() error {
	log.Printf("Set up GRPC server\n")

	// Validate address
	addr, err := net.ResolveTCPAddr("tcp", a.Config.GRPCAddr)
	if err != nil {
		return err
	}

	if addr == nil {
		return errors.New("no address provided")
	}

	var srv Server

	log.Printf("GRPC server listens on %s\n", a.Config.GRPCAddr)
	listener, err := net.Listen("tcp", a.Config.GRPCAddr)
	if err != nil {
		return err
	}
	srv.listener = listener

	kp := keepalive.ServerParameters{MaxConnectionAge: time.Minute}

	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(kp),
	}
	srv.grpcServer = grpc.NewServer(opts...)

	log.Printf("Create library repo client connection to GRPC server %s\n", a.Config.LibraryRepoAddr)
	libraryRepoClient, err := libraryrepoClientDev.New(a.Config.LibraryRepoAddr, a.Config.DBConnectString, a.Config.DBName)
	if err != nil {
		return errors.New(err.Error())
	}

	log.Printf("Create library repo client writer connection\n")
	writerConn, err := makeWriterConn(a.Config.LibraryRepoAddr)
	if err != nil {
		panic(errors.New(err.Error()))
	}

	log.Printf("Create library server using library repo client writer connection\n")
	librarySrv := library_server.New(libraryRepoClient, writerConn)
	librarypb.RegisterLibraryServer(srv.GRPCServer(), librarySrv)

	a.server = &srv

	return nil
}

func makeWriterConn(addr string) (*grpc.ClientConn, error) {
	if addr == "" {
		return nil, nil
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
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

	r := mux.NewRouter()
	r.HandleFunc("/health", makeHealthCheckHandler()).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:    addr.String(),
		Handler: r,
	}
	srv.Addr = a.Config.HTTPAddr

	a.server.httpServer = srv

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

// GRPCServer returns the server's grpc.Server.
func (srv *Server) GRPCServer() *grpc.Server {
	return srv.grpcServer
}

// Stop stops the gRPC server.
func (srv *Server) Stop() {
	srv.grpcServer.GracefulStop()
	srv.httpServer.Shutdown(context.TODO())
}

// ServeForever listens for gRPC requests.
func (srv *Server) ServeForever() error {
	// Serve GRPC server
	log.Printf("Starting up GRPC server\n")
	errServerCh := make(chan error)
	go func() {
		err := srv.grpcServer.Serve(srv.listener)
		if err != nil {
			errServerCh <- err
		}
		errServerCh <- nil
	}()

	// Serve HTTP server
	log.Printf("Starting up HTTP server\n")
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
