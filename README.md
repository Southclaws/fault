# fault

> _because the word error is overused_

Fault is a simple and minimal error utility library designed to help you diagnose problems in your application logic without overly verbose stack traces and by allowing you to annotate wrapped errors with metadata that plays nicely with structured logging tools and JSON as well as the existing Go errors ecosystem.

It's not at a v1 state yet, the API will likely change as I refine the ergonomics based on feedback.

## Example

```go
func (r *Repository) GetUser(ctx context.Context, userID string) (*User, error) {
    user, err := r.db.Query(`... etc`)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fault.WithValue(err, "user not found", "user_id", userID)
        }
        return nil, fault.WithValue(err, "failed to execute query", "user_id", userID)
    }

    return user, nil
}

// In the HTTP handler...

func (u *usersController) getUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "user_id")

    user, err := u.repository.GetUser(r.Context(), id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            w.WriteHeader(404)
            return
        }

        // Simple, easy, useful structured logging from errors!
        log.Error("failed to get user", fault.Context(err))

        w.WriteHeader(500)
        return
    }

    // ... etc
}

// When this fails, you'll see this nice and simple structured log entry in your logs:
// {
//   "message": "failed to execute query: sql database fell asleep",
//   "values": {
//     "user_id": "xyz123"
//   },
//   "trace": [
//     {
//       "message": "failed to execute query: sql database fell asleep",
//       "location": "users/repository.go:14"
//     },
//     {
//       "message": "sql database fell asleep",
//       "location": "users/db.go:22"
//     }
//   ]
// }
```

Fault takes a lot of inspiration from the [context](https://pkg.go.dev/context) package. Implements a similar mechanism to `context.WithValue` where you can annotate an element in the chain with some arbitrary string-based metadata.

You can use this alongside pkg/errors and pretty much any other error library. It's all based on standard and commonly used interfaces.

## Usage

### `fault.New`

Simply creates a new error, just like all the other libraries. With the added benefit of storing some minimal source code location information.

### `fault.WithValue`

Wraps an error with a piece of key-value data. Can later be extracted as structured data. Also stores source code location.

### `fault.WithMessage`

Wraps an error with an additional message. Same as `Wrap` in a lot of other error libraries.

### `fault.Context`

Returns a data structure representing an error chain, source code locations as well as any contextual values contributed to the error chain. This is ideal for using with structured logging.

### JSON

This library also contains a simple JSON encoding implementation which uses the result of `fault.Context` to encode an error object into a structured piece of JSON for logging or HTTP responses.

## Why?

The motivation for this package came about due to frustration with the disparity between the simple non-structured, string-based errors of Go and the useful, searchability of structured logging. There wasn't a simple glue utility to bridge the gap between dealing with an error in a service layer or HTTP request handler and the structured logging solution. On a small team without a bespoke solution, this can lead to error messages being the only way to diagnose problems. Simple error strings are great, but they lack the ergonomics to easily embed structured metadata.

### Error Context

The canonical standard library approach is:

```go
var ErrSentinelError = errors.New("something bad happened")

func PublicAPI(ctx context.Context, userID string, transactionID string, paymentReference string) (*Thing, error) {
    if badThing {
        return nil, ErrSentinelError
    }
}
```

In this example, this simple sentinel error makes it easy to find the _location_ of the problem - you just search for the error message in the codebase and you find this file. The problem is, the arguments that will help diagnose what went wrong are lost. `userID`, `transactionID` and `paymentReference` are all lost.

So, engineers often make use of error wrapping, with Dave Cheney's errors package (or one of the many similar libraries) like so:

```go
var ErrSentinelError = errors.New("something bad happened")

func PublicAPI(ctx context.Context, userID string, transactionID string, paymentReference string) (*Thing, error) {
    if badThing {
        return nil, errors.Wrapf(ErrSentinelError,
            "failed to do the thing with user %s on transaction %s with payment reference %s",
            userID, transactionID, paymentReference)
    }
}
```

Now you can see the user's ID, transaction ID and payment reference in the error message. Which is great, if there's a mistake in the business logic, you have more information to help reproduce and diagnose. Maybe it's a problem specific to the user or transaction and having that information is vital to solving issues fast.

But this is still not great. It lacks ergonomics, and when APIs aren't ergonomic, they don't get used. In a lot of the codebases I've worked on, I rarely see engineers use this pattern to decorate errors. Sure, you can try to enforce this during review, but these are all opportunities to forget still. Ergonomic APIs make developers happier because it makes their job easier. And we're all writing code to make tasks easier, not harder.

```
"failed to do the thing with user %s on transaction %s with payment reference %s"
"failed to do the thing with user ID %s on transaction ID %s with payment reference %s"
"failed to do the thing: user \"%s\" transaction \"%s\" payment reference \"%s\""
"failed to do the thing: user [%s], transaction [%s], payment reference [%s]"
"failed to do the thing: user:%s transaction:%s payment_reference:%s"
```

Another issue is how this is completely unstructured. This is an English sentence with three pieces of _data_ sprinkled within it. While working on a large, continually evolving, codebase with a team of engineers, consistency can sometimes get forgotten. Varying ways of writing errors results in difficult log searches. Structured logging was created to resolve this, but there's no way to turn the above error into a properly structured log entry.

```json
failed to do the thing {"user_id":"1234","transaction_id":"9876","payment_reference":"xyzabc"}
```

```logfmt
msg="failed to do the thing" user_id="1234" transaction_id="9876" payment_reference="xyzabc"
```

```json
{
    "message": "failed to do the thing",
    "user_id": "1234",
    "transaction_id": "9876",
    "payment_reference": "xyzabc"
}
```

_this kind of logging makes searching for ALL the logs events for a specific user trivial_

And finally, structured logging and consistent naming of keys makes security a lot easier. You can configure your log aggregator to strip out or limit the exposure of `paymentReference` values if you don't want those in logs. This is not possible with basic string-based errors.

```json
{
    "message": "failed to do the thing",
    "user_id": "1234",
    "transaction_id": "9876",
    "payment_reference": "xyz[REDACTED]"
}
```

### Stack Traces

(Hot take alert!) Stack traces are 80% useless.

Here's a simple trace from an error using the popular pkg/errors library:

```
=== RUN   TestGiveMeAnError
oh no
github.com/Southclaws/fault.SomethingBadHappened
 /Users/southclaws/Work/fault/fault_test.go:12
github.com/Southclaws/fault.WouldBeAShameIfSomethingBadHappened
 /Users/southclaws/Work/fault/fault_test.go:16
github.com/Southclaws/fault.IReallyHopeNothingBadHappens
 /Users/southclaws/Work/fault/fault_test.go:20
github.com/Southclaws/fault.NothingBadWillHappen
 /Users/southclaws/Work/fault/fault_test.go:24
github.com/Southclaws/fault.TestGiveMeAnError
 /Users/southclaws/Work/fault/fault_test.go:28
testing.tRunner
 /usr/local/go/src/testing/testing.go:1439
runtime.goexit
```

And here's what it looks like when printed using the popular structured logging library Zap:

```
/usr/local/go/src/runtime/asm_arm64.s:12632022-07-10T18:30:29.222+0100 ERROR fault/fault_test.go:34 i failed you {"error": "oh no", "errorVerbose": "oh no\ngithub.com/Southclaws/fault.SomethingBadHappened\n\t/Users/southclaws/Work/fault/fault_test.go:12\ngithub.com/Southclaws/fault.WouldBeAShameIfSomethingBadHappened\n\t/Users/southclaws/Work/fault/fault_test.go:16\ngithub.com/Southclaws/fault.IReallyHopeNothingBadHappens\n\t/Users/southclaws/Work/fault/fault_test.go:20\ngithub.com/Southclaws/fault.NothingBadWillHappen\n\t/Users/southclaws/Work/fault/fault_test.go:24\ngithub.com/Southclaws/fault.TestGiveMeAnError\n\t/Users/southclaws/Work/fault/fault_test.go:28\ntesting.tRunner\n\t/usr/local/go/src/testing/testing.go:1439\nruntime.goexit\n\t/usr/local/go/src/runtime/asm_arm64.s:1263"}
github.com/Southclaws/fault.TestGiveMeAnError
 /Users/southclaws/Work/fault/fault_test.go:34
testing.tRunner
 /usr/local/go/src/testing/testing.go:1439
```

And Zero:

```
{"level":"error","stack":[{"func":"SomethingBadHappened","line":"15","source":"fault_test.go"},{"func":"WouldBeAShameIfSomethingBadHappened","line":"19","source":"fault_test.go"},{"func":"IReallyHopeNothingBadHappens","line":"23","source":"fault_test.go"},{"func":"NothingBadWillHappen","line":"27","source":"fault_test.go"},{"func":"TestGiveMeAnError","line":"31","source":"fault_test.go"},{"func":"tRunner","line":"1439","source":"testing.go"},{"func":"goexit","line":"1263","source":"asm_arm64.s"}],"error":"oh no","time":"2022-07-10T18:33:45+01:00"}
```

There's a lot of noise here. Zero is the best as it actually structures the trace properly but it's still including some unnecessary elements like `asm_arm64.s` (assembly in the Go runtime) and `goexit`.

Most of the time, your errors will be directly related to business logic if you're building a product. So all this low level noise is useless. Especially when it's formatted as a single line with a ton of `\n` characters. Ideally, you want something that's going to play well with your log aggregator with zero configuration.

### Error Categories

This is something I'd like to solve using this library, but I want to figure out a good API first. So this is more of a roadmap item currently. Though you can implement this with the current library already.

All problems in a typical API server can be categorised. We do this already using HTTP statuses. If you try to load a user and it doesn't exist, that's a "404 Not found" category of error. If you try to write data to the database and it crashes, that's a "500 Internal server error". The HTTP specification has figured out all the possible kinds of problem and encoded them into status codes. We also have the same with gRPC, there are 17 different types of error documented [here](https://grpc.github.io/grpc/core/md_doc_statuscodes.html) and any problem your code runs into can be expressed using one of these categories.

The problem with simple Go errors is you end up handling these explicitly using `errors.Is` in HTTP handlers. This is fine for small applications, it's simple, elegant and really readable without any abstractions. But large applications, that may be using OpenAPI or protocol buffers need a way to pragmatically map application errors to on-the-wire error types.

This is also a matter of the separation of concerns. Should your HTTP handler be _explicitly_ checking for a `sql.ErrNoRows` error? That's a leaky abstraction. If you use any form of domain driven design or MVC, you know that the layers should not expose internal details to themselves. An error specific to SQL being exposed all the way up at the interface layer is a bad design.

So error categories is something I'd like to resolve to make my own products and teams work more efficiently. There are some solutions out there already but they are quite complex and specific to the project.
