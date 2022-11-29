package tests

import (
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
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

func externalWrappedPostgresError() error {
	err := pgerror()
	return fmt.Errorf("external pg error: %w", err)
}

func pgerror() error {
	return &pgconn.PgError{
		Severity:         "fatal",
		Code:             "123",
		Message:          "your sql was wrong bro",
		Detail:           "very wrong",
		Hint:             "rtfm",
		Position:         420,
		InternalPosition: 69,
		InternalQuery:    "internal_query",
		Where:            "where",
		SchemaName:       "schema_name",
		TableName:        "table_name",
		ColumnName:       "column_name",
		DataTypeName:     "data_type_name",
		ConstraintName:   "constraint_name",
		File:             "file",
		Line:             0,
		Routine:          "routine",
	}
}
