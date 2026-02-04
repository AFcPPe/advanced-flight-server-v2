package errs

import "errors"

// ErrNotAuthenticated 首包不是$ID，需要断开连接
var ErrNotAuthenticated = errors.New("first packet must be $ID")
