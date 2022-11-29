package tests

import (
	"errors"
	"fmt"

	pkg_errors "github.com/pkg/errors"
)

var (
	errExternalStdlib = errors.New("stdlib external error")
	errExternalPkgerr = pkg_errors.New("github.com/pkg/errors external error")
)

func externalError() error {
	return fmt.Errorf("external error wrapped with errorf: %w", errExternalStdlib)
}

func externalWrappedError() error {
	return pkg_errors.Wrap(errExternalPkgerr, "external error wrapped with pkg/errors")
}
