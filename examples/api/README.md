# API Example

This is an example project to demonstrate the use of Fault in an API context. This fake API expose a single endpoint 
`/users/{id}`. Asking for user 999 will generate an error by design. Asking for 123 will return a user info payload.
Asking for anything else will return a "NotFound" response.

## Running

From the repository root, using golang 1.19+, run

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
2023-06-19T00:41:29.662-0400	ERROR	http/http.go:105
	db error: connection lost
		examples/api/main.go:58
	Could not get user
		examples/api/main.go:59
	{"http_method": "GET", "request_id": "VRyQ76vHEy-000001", "request_path": "/users/999", "remote_ip": "127.0.0.1:51715", "protocol": "HTTP/1.1", "user_id": "999", "error": "Could not get user: db error: connection lost"}
2023-06-19T00:41:29.662-0400	INFO	http/http.go:64	API Request	{"protocol": "HTTP/1.1", "http_method": "GET", "request_id": "VRyQ76vHEy-000001", "request_path": "/users/999", "remote_ip": "127.0.0.1:51715", "status": 500, "latency": "130.682µs"}
2023-06-19T00:41:34.417-0400	ERROR	http/http.go:105
	db error: user id[321] not found
		examples/api/main.go:67
	User not found
		examples/api/main.go:68
	{"user_id": "321", "remote_ip": "127.0.0.1:51716", "protocol": "HTTP/1.1", "request_id": "VRyQ76vHEy-000002", "http_method": "GET", "request_path": "/users/321", "error": "User not found: db error: user id[321] not found"}
2023-06-19T00:41:34.417-0400	INFO	http/http.go:64	API Request	{"request_id": "VRyQ76vHEy-000002", "http_method": "GET", "request_path": "/users/321", "remote_ip": "127.0.0.1:51716", "protocol": "HTTP/1.1", "status": 404, "latency": "102.541µs"}
2023-06-19T00:42:35.202-0400	INFO	http/http.go:64	API Request	{"http_method": "GET", "request_id": "VRyQ76vHEy-000003", "request_path": "/users/123", "remote_ip": "127.0.0.1:51724", "protocol": "HTTP/1.1", "status": 200, "latency": "107.706µs"}
```

## How

There is a couple of things we're trying to demonstrate here:

1. We use fault's flatten method to produce the equivalent of a stacktrace that we display with the error
2. We use fault's `ftag` to "tag" the errors and infer the http status code from that
3. We use fault's `fctx` to add http context fields, but also controller based values, like in this example, the `userId`
4. We use fault's user-friendly error messages while building the user-facing message we return as part of the http response. (check the curl responses)