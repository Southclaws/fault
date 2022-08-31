# fault

> _because the word error is overused_

Fault is a simple and minimal collection of error utilities designed to help you diagnose problems in your application logic without the need for overly verbose stack traces and by allowing you to annotate wrapped errors with metadata that plays nicely with structured logging tools and JSON as well as the existing Go errors ecosystem.

It's not at a v1 state yet, the API will likely change as I refine the ergonomics based on feedback.

## Libraries

- errctx: bridging the gap between `context.Context` and `error`.
- errmeta: structured simple string-based key-value metadata in errors.
- errtag: well defined error kinds such as "not found" or "validation failure".

Documentation: https://pkg.go.dev/github.com/Southclaws/fault

Inspired by:

- https://github.com/cockroachdb/logtags
- https://github.com/cockroachdb/errors/tree/master/contexttags
- https://pkg.go.dev/google.golang.org/grpc/status
