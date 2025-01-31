package idgen

import "github.com/rs/xid"

func New[T ~string]() T {
	return T(xid.New().String())
}
