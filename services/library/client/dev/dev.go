package dev

import (
	"library-under-the-sea/services/library/client/grpc"
	"library-under-the-sea/services/library/client/logical"
	library "library-under-the-sea/services/library/domain"
)

// New returns a gRPC client if the gRPC address is set via a flag. If not,
// it returns a logical client.
func New() (library.Client, error) {
	if grpc.IsEnabled() {
		return grpc.New()
	}

	return logical.New(r), nil
}
