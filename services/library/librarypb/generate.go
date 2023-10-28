package librarypb

//go:generate protoc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false --go_out=. --go-grpc_out=. -I=. ./librarypb.proto
