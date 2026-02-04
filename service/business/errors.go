package business

import (
	"fmt"

	"connectrpc.com/connect"
)

var (
	ErrorUnspecifiedID      = connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("no id was supplied"))
	ErrorEmptyValueSupplied = connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("empty value supplied"))
	ErrorItemExist          = connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("specified item already exists"))
	ErrorItemDoesNotExist   = connect.NewError(connect.CodeNotFound, fmt.Errorf("specified item does not exist"))
	ErrorInitializationFail = connect.NewError(connect.CodeInternal, fmt.Errorf("internal configuration is invalid"))
)
