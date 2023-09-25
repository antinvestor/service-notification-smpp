package business

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrorUnspecifiedID      = status.Error(codes.InvalidArgument, "No id was supplied")
	ErrorEmptyValueSupplied = status.Error(codes.InvalidArgument, "Empty value supplied")
	ErrorItemExist          = status.Error(codes.AlreadyExists, "Specified item already exists")
	ErrorItemDoesNotExist   = status.Error(codes.NotFound, "Specified item does not exist")

	ErrorInitializationFail = status.Error(codes.Internal, "Internal configuration is invalid")
)
