package errlib

import "github.com/cockroachdb/errors"

func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}
