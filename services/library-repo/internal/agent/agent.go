package agent

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"library-under-the-sea/services/library-repo/libraryrepopb"
	libraryrepo_server "library-under-the-sea/services/library-repo/server"
)

type Config struct {
	// HTTPAddr is the healthcheck address
	HTTPAddr string
	// GRPCAddr is the GRPC address
	GRPCAddr string
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

	listener, err := net.Listen("tcp", a.Config.GRPCAddr)
	if err != nil {
		return err
	}
	srv.listener = listener

	kp := keepalive.ServerParameters{MaxConnectionAge: time.Minute}

	opts := []grpc.ServerOption{grpc.KeepaliveParams(kp)}
	srv.grpcServer = grpc.NewServer(opts...)

	libraryrepoSrv := libraryrepo_server.New(a.Config.DBConnectString, a.Config.DBName)
	libraryrepopb.RegisterLibraryRepoServer(srv.GRPCServer(), libraryrepoSrv)

	a.server = &srv

	return nil
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
