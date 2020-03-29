package errs

import (
"fmt"
)

var (
	ErrBadRequest = fmt.Errorf("bad request")
	ErrUnauthorized = fmt.Errorf("unauthorized")

	ErrNoRows = fmt.Errorf("sql: no rows")
)

