package utils

const (
	ErrDataNotFound          = "sql: no rows in result set"
	ErrIDCannotBeEmpty       = "ID cannot be empty"
	ErrUsernameAlreadyExists = "username already exists"
	ErrPhoneAlreadyExists    = "phone already exists"
	ErrEmailAlreadyExists    = "email already exists"
	ErrPasswordNotMatch      = "password not match"
	ErrInvalidPassword       = "invalid password"
)

const (
	DateLayout     = "2006-01-02"
	DatetimeLayout = "2006-01-02 15:04:05"
)

const (
	InvalidTokenMsg     = "invalid token data"
	NotAllowedModifyMsg = "user not allowed to modify this data"
	NotAllowedAccessMsg = "user not allowed to access this data"
)

var (
	SortList = []string{"ASC", "DESC"}
)

// GRPC STATUS CODE

const (
	OK                  = 0
	CANCELLED           = 1
	UNKNOWN             = 2
	INVALID_ARGUMENT    = 3
	DEADLINE_EXCEEDED   = 4
	NOT_FOUND           = 5
	ALREADY_EXISTS      = 6
	PERMISSION_DENIED   = 7
	RESOURCE_EXHAUSTED  = 8
	FAILED_PRECONDITION = 9
	ABORTED             = 10
	OUT_OF_RANGE        = 11
	UNIMPLEMENTED       = 12
	INTERNAL            = 13
	UNAVAILABLE         = 14
	DATA_LOSS           = 15
	UNAUTHENTICATED     = 16
)

const (
	NoRowsInResultSet = "no rows in result set"
	InvalidEventType  = "invalid event type"
)

const (
	RolePicOutsourcing = "pic_outsourcing" //outsource
	RoleSuperAdmin     = "admin"           //admin
)
