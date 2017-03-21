package lfstransfers3

import "errors"

var (
	EINVALIDDIRECTORY              = errors.New("Invalid Directory")
	EPART_SIZE_LARGER_THAN_ALLOWED = errors.New("")
)
