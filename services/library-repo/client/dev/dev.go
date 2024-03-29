package dev

import (
	"library-under-the-sea/services/library-repo/client/grpc"
	"library-under-the-sea/services/library-repo/client/logical"
	library_repo "library-under-the-sea/services/library-repo/domain"
)

// New returns a gRPC client if the gRPC address is set via a flag. If not,
// it returns a logical client.
func New(addr string, connectString string, dbName string) (library_repo.Client, error) {
	if addr != "" {
		return grpc.New(addr)
	}

	return logical.New(connectString, dbName), nil
}
