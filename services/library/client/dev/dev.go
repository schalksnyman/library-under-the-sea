package dev

import (
	libraryrepo "library-under-the-sea/services/library-repo/domain"
	"library-under-the-sea/services/library/client/grpc"
	"library-under-the-sea/services/library/client/logical"
	library "library-under-the-sea/services/library/domain"
)

// New returns a gRPC client if the gRPC address is set via a flag. If not,
// it returns a logical client.
func New(r libraryrepo.Client) (library.Client, error) {
	if grpc.IsEnabled() {
		return grpc.New()
	}

	return logical.New(r), nil
}
