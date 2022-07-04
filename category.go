package fault

// Kind represents a type of problem. These are intended to be broadly applying
// to common types of error thrown any application. They are based on the gRPC
// status codes. The reason for this is the list is quite small and pretty much
// any problem that occurs in an application will fit into one of these types.
//
// For more info, see: https://grpc.github.io/grpc/core/md_doc_statuscodes.html
//
type Kind uint8

const (
	// Not an error; returned on success.
	Ok Kind = 0

	// The operation was cancelled, typically by the caller.
	Cancelled Kind = 1

	// Unknown error. For example, this error may be returned when a Status value received from another address space belongs to an error space that is not known in this address space. Also errors raised by APIs that do not return enough error information may be converted to this error.
	Unknown Kind = 2

	// The client specified an invalid argument. Note that this differs from FAILED_PRECONDITION. INVALID_ARGUMENT indicates arguments that are problematic regardless of the state of the system (e.g., a malformed file name).
	InvalidArgument Kind = 3

	// The deadline expired before the operation could complete. For operations that change the state of the system, this error may be returned even if the operation has completed successfully. For example, a successful response from a server could have been delayed long
	DeadlineExceeded Kind = 4

	// Some requested entity (e.g., file or directory) was not found. Note to server developers: if a request is denied for an entire class of users, such as gradual feature rollout or undocumented allowlist, NOT_FOUND may be used. If a request is denied for some users within a class of users, such as user-based access control, PERMISSION_DENIED must be used.
	NotFound Kind = 5

	// The entity that a client attempted to create (e.g., file or directory) already exists.
	AlreadyExists Kind = 6

	// The caller does not have permission to execute the specified operation. PERMISSION_DENIED must not be used for rejections caused by exhausting some resource (use RESOURCE_EXHAUSTED instead for those errors). PERMISSION_DENIED must not be used if the caller can not be identified (use UNAUTHENTICATED instead for those errors). This error code does not imply the request is valid or the requested entity exists or satisfies other pre-conditions.
	PermissionDenied Kind = 7

	// Some resource has been exhausted, perhaps a per-user quota, or perhaps the entire file system is out of space.
	ResourceExhausted Kind = 8

	// The operation was rejected because the system is not in a state required for the operation's execution. For example, the directory to be deleted is non-empty, an rmdir operation is applied to a non-directory, etc. Service implementors can use the following guidelines to decide between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE: (a) Use UNAVAILABLE if the client can retry just the failing call. (b) Use ABORTED if the client should retry at a higher level (e.g., when a client-specified test-and-set fails, indicating the client should restart a read-modify-write sequence). (c) Use FAILED_PRECONDITION if the client should not retry until the system state has been explicitly fixed. E.g., if an "rmdir" fails because the directory is non-empty, FAILED_PRECONDITION should be returned since the client should not retry unless the files are deleted from the directory.
	FailedPrecondition Kind = 9

	// The operation was aborted, typically due to a concurrency issue such as a sequencer check failure or transaction abort. See the guidelines above for deciding between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE.
	Aborted Kind = 10

	// The operation was attempted past the valid range. E.g., seeking or reading past end-of-file. Unlike INVALID_ARGUMENT, this error indicates a problem that may be fixed if the system state changes. For example, a 32-bit file system will generate INVALID_ARGUMENT if asked to read at an offset that is not in the range [0,2^32-1], but it will generate OUT_OF_RANGE if asked to read from an offset past the current file size. There is a fair bit of overlap between FAILED_PRECONDITION and OUT_OF_RANGE. We recommend using OUT_OF_RANGE (the more specific error) when it applies so that callers who are iterating through a space can easily look for an OUT_OF_RANGE error to detect when they are done.
	OutOfRange Kind = 11

	// The operation is not implemented or is not supported/enabled in this service.
	Unimplemented Kind = 12

	// Internal errors. This means that some invariants expected by the underlying system have been broken. This error code is reserved for serious errors.
	Internal Kind = 13

	// The service is currently unavailable. This is most likely a transient condition, which can be corrected by retrying with a backoff. Note that it is not always safe to retry non-idempotent operations.
	Unavailable Kind = 14

	// Unrecoverable data loss or corruption.
	DataLoss Kind = 15

	// The request does not have valid authentication credentials for the operation.
	Unauthenticated Kind = 16
)

func (k Kind) String() string {
	switch k {
	case Ok:
		return "OK"
	case Cancelled:
		return "CANCELLED"
	case Unknown:
		return "UNKNOWN"
	case InvalidArgument:
		return "INVALID_ARGUMENT"
	case DeadlineExceeded:
		return "DEADLINE_EXCEEDED"
	case NotFound:
		return "NOT_FOUND"
	case AlreadyExists:
		return "ALREADY_EXISTS"
	case PermissionDenied:
		return "PERMISSION_DENIED"
	case ResourceExhausted:
		return "RESOURCE_EXHAUSTED"
	case FailedPrecondition:
		return "FAILED_PRECONDITION"
	case Aborted:
		return "ABORTED"
	case OutOfRange:
		return "OUT_OF_RANGE"
	case Unimplemented:
		return "UNIMPLEMENTED"
	case Internal:
		return "INTERNAL"
	case Unavailable:
		return "UNAVAILABLE"
	case DataLoss:
		return "DATA_LOSS"
	case Unauthenticated:
		return "UNAUTHENTICATED"
	}
	return "UNKNOWN"
}

// func (k Kind) Status() int {
// 	switch k {
// 	case Ok:
// 		return http.StatusOK
// 	case Cancelled:
// 		return http.StatusGatewayTimeout
// 	case Unknown:
// 		return "UNKNOWN"
// 	case InvalidArgument:
// 		return "INVALID_ARGUMENT"
// 	case DeadlineExceeded:
// 		return "DEADLINE_EXCEEDED"
// 	case NotFound:
// 		return "NOT_FOUND"
// 	case AlreadyExists:
// 		return "ALREADY_EXISTS"
// 	case PermissionDenied:
// 		return "PERMISSION_DENIED"
// 	case ResourceExhausted:
// 		return "RESOURCE_EXHAUSTED"
// 	case FailedPrecondition:
// 		return "FAILED_PRECONDITION"
// 	case Aborted:
// 		return "ABORTED"
// 	case OutOfRange:
// 		return "OUT_OF_RANGE"
// 	case Unimplemented:
// 		return "UNIMPLEMENTED"
// 	case Internal:
// 		return "INTERNAL"
// 	case Unavailable:
// 		return "UNAVAILABLE"
// 	case DataLoss:
// 		return "DATA_LOSS"
// 	case Unauthenticated:
// 		return "UNAUTHENTICATED"
// 	}
// 	return "UNKNOWN"
// }
