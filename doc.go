// Package fault provides a mechanism for wrapping errors with various types of
// useful metadata. It implements this as a kind of middleware style pattern by
// providing a simple interface that can be passed to a call to `fault.Wrap`.
//
// # Rationale
//
// The reason for this is because nesting a lot of calls to various `Wrap` calls
// is really awkward to write and read. The Golang errors ecosystem is diverse
// but unfortunately, composing together many small error related tools remains
// awkward due to the simple yet difficult to extend patterns set by the Golang
// standard library and popular error packages.
//
// For example, to combine pkg/errors, tracerr and errctx you'd have to write:
//
//	errctx.Wrap(errors.Wrap(tracerr.Wrap(err), "failed to get user"), ctx)
//
// Which is a bit of a nightmare to write (many nested calls) and a nightmare to
// read (not clear where the arguments start and end for each function). Because
// of this, it's not common to compose together libraries from the ecosystem.
//
// # Solution
//
// To resolve that, Fault provides a simpler and more ergonomic way to decorate
// an error chain with additional information at any point during a call stack.
// Internally, it works like any error wrapping library and satisfies all the
// common interfaces. The `Wrap` function, however, takes a variadic list of
// what looks like, if you squint, middleware functions!
//
// These utilities all satisfy a common interface so anyone can write a package
// which works with Fault and provides some kind of additional information for
// an error message and Fault handles the rest.
//
// # Usage: Wrapping
//
// Wrapping errors as you ascend the call stack is essential to providing your
// team with adequate context when something goes wrong. This is akin to trace
// spans but in reverse order as it's dealing with functions at return-time
// rather than call-time.
//
// Simply wrap your errors as you would with any library and pass any built-in
// modifiers that you need.
//
//	if err != nil {
//		return fault.Wrap(err,
//			errctx.Ctx(ctx),
//			errmeta.Field("user_id", userID),
//			issues.Issue("A message intended for the end-user to read."),
//		)
//	}
//
// Fault provides a small collection of built-in decorators which cover various
// use-cases for most medium sized applications.
//
// # Usage: Handling
//
// Wrapping errors is only half the story, eventually you'll need to actually
// *handle* the error (and no, `return err` is not "handling" an error, it's
// saying "I don't know what to do! Caller, you deal with this!".)
//
// Handling a Fault error means you need to get all the useful structured data
// back out of the error so you can do something useful with it. This is usually
// (but not limited to)
//
// 1. Log with your favourite structured logging tool
// 2. Construct an API response to send to the client, describing what happened
// 3. Interface with your altering, tracing, monitoring system
//
// All of these kinds of tools require some kind of structure, and unfortunately
// `err.Error()` just doesn't cut it.
//
// Depending on what middleware you've used, you'll be using different extractor
// functions to pull out the structured data so check the corresponding docs.
package fault
