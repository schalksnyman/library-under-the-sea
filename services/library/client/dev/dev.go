package dev

import (
	libraryRepo "library-under-the-sea/services/library-repo/domain"
	"library-under-the-sea/services/library/client/grpc"
	"library-under-the-sea/services/library/client/logical"
	library "library-under-the-sea/services/library/domain"
)

// New returns a gRPC client if the gRPC address is set via a flag. If not,
// it returns a logical client.
func New(addr string, r libraryRepo.Client) (library.Client, error) {
	if addr != "" {
		return grpc.New(addr)
	}

	return logical.New(r), nil
}
