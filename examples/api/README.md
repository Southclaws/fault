# API Example

This is an example project to demonstrate the use of Fault in an API context. This fake API expose a single endpoint 
`/users/{id}`. Asking for user 999 will generate an error by design. Asking for 123 will return a user info payload.
Asking for anything else will return a "NotFound" response.

## Running

From the repository root, using golang 1.21+, run (because the slog package is needed).

```bash
go run ./examples/api/main.go
```

Then you can run the following curl commands
```bash
curl http://localhost:3333/users/999 # wil return a 500 error
curl http://localhost:3333/users/321 # will return not found
curl http://localhost:3333/users/123 # will return a user in JSON
```

Running each curl should produce the following logs
```
2023/09/09 11:05:06 ERROR 
db error: connection lost
	examples/api/main.go:53
Could not get user
	examples/api/main.go:54
<ftag>
	examples/api/main.go:54
 http_method=GET user_id=999 request_id=It73FDo3WC-000005 request_path=/users/999 remote_ip=127.0.0.1:52246 protocol=HTTP/1.1 error="Could not get user: db error: connection lost"
2023/09/09 11:05:06 INFO API Request request_id=It73FDo3WC-000005 request_path=/users/999 remote_ip=127.0.0.1:52246 protocol=HTTP/1.1 http_method=GET status=500 latency=96.25µs

2023/09/09 11:05:25 ERROR 
db error: user id[321] not found
    examples/api/main.go:62
User not found
    examples/api/main.go:63
<ftag>
    examples/api/main.go:63
 http_method=GET user_id=321 request_path=/users/321 remote_ip=127.0.0.1:52246 request_id=It73FDo3WC-000006 protocol=HTTP/1.1 error="User not found: db error: user id[321] not found"
2023/09/09 11:05:25 INFO API Request protocol=HTTP/1.1 http_method=GET request_path=/users/321 remote_ip=127.0.0.1:52246 request_id=It73FDo3WC-000006 status=404 latency=65.306µs

2023/09/09 11:06:28 INFO API Request protocol=HTTP/1.1 request_id=It73FDo3WC-000008 http_method=GET request_path=/users/123 remote_ip=127.0.0.1:52426 status=200 latency=46.562µs

```

## How

There is a couple of things we're trying to demonstrate here:

1. We use fault's flatten method to produce the equivalent of a stacktrace that we display with the error
2. We use fault's `ftag` to "tag" the errors and infer the http status code from that
3. We use fault's `fctx` to add http context fields, but also controller based values, like in this example, the `userId`
4. We use fault's user-friendly error messages while building the user-facing message we return as part of the http response. (check the curl responses)