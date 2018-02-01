package auth

import "errors"

var ErrNotFound = errors.New("key not found")

type CredentialsStore interface {
	Set(key, value string)
	Get(key string, clb func(value string)) error
}
