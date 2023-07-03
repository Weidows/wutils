package log

import (
	"github.com/morrisxyang/errors"
	"os"
)

func a() error {
	err := b()
	err = errors.Wrap(err, "a failed reason")
	return err
}

func b() error {
	err := c()
	err = errors.Wrap(err, "b failed reason")
	return err
}

func c() error {
	_, err := os.Open("test")
	if err != nil {
		return errors.WrapWithCode(err, 123, "c failed reason")
	}
	return nil
}
