package tests

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
)

func errorCaller(kind int) error {
	err := errorCallerFromMiddleOfChain(kind)
	if err != nil {
		return fault.Wrap(err)
	}

	return nil
}

func errorCallerFromMiddleOfChain(kind int) error {
	err := errorProducerFromRootCause(kind)
	if err != nil {
		return fault.Wrap(err, fmsg.With("failed to call function"))
	}

	return nil
}

func errorProducerFromRootCause(kind int) error {
	err := rootCause(kind)
	if err != nil {
		return fault.Wrap(err)
	}

	return nil
}
