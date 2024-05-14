package constant

type LogLevelType string

const (
	Info  LogLevelType = "INFO"
	Warn  LogLevelType = "WARN"
	Error LogLevelType = "ERROR"
	Debug LogLevelType = "DEBUG"
)

type ErrorEnums struct {
	ErrorCode    int
	ErrorMessage string
}

const BaseApiPath = "/fast_storage/api/v1"

const LogFileLocation = "fast_storage_service_log_%d_%d_%d.log"
const DeltaPositive = 0.5
const DeltaNegative = -0.5
const YyyyMmDdHhMmSsFormat = "2006-01-02 15:04:05"
const AscKeyword = "ASC"
const DescKeyword = "DESC"
const EmptyString = ""
const (
	ContentTypeBinary    = "application/octet-stream"
	ContentTypeForm      = "application/x-www-form-urlencoded"
	ContentTypeJSON      = "application/json"
	ContentTypeHTML      = "text/html; charset=utf-8"
	ContentTypeText      = "text/plain; charset=utf-8"
	ContentTypeIconImage = "image/x-icon"
	ContentTypePngImage  = "image/png"
)

const AuthenticationFailed = "AUTHENTICATION_FAILED"
const AuthenticationCorrupted = "AUTHENTICATION_CORRUPTED"
const AuthenticationSuccessfully = "AUTHENTICATION_SUCCESSFULLY"

var SensitiveField = [...]string{"password", "jwt", "token", "client_secret", "Authorization", "x-api-key"} // [...] instead of []: it ensures you get a (fixed size) array instead of a slice. So the values aren't fixed but the size is.
var ValidMethod = []string{"GET", "POST", "PUT", "DELETE"}

type LogKey string

const UsernameLogKey LogKey = "username"
const UserIdLogKey LogKey = "userId"
const TraceIdLogKey LogKey = "traceId"
const LogPattern = "[%s] [%s] [%s] 👉️ \t%s"

var (
	Success = ErrorEnums{
		ErrorCode:    0,
		ErrorMessage: "Success",
	}
	InternalFailure = ErrorEnums{
		ErrorCode:    -1,
		ErrorMessage: "An error has been occurred, please try again later",
	}
	PageNotFound = ErrorEnums{
		ErrorCode:    -2,
		ErrorMessage: "You're consuming an unknow endpoint, please check your url (404 Page Not Found)",
	}
	MethodNotAllowed = ErrorEnums{
		ErrorCode:    -3,
		ErrorMessage: "This url is configured method that not match with your current method, please check again (405 Method Not Allowed)",
	}
	QueryStatementError = ErrorEnums{
		ErrorCode:    -4,
		ErrorMessage: "Query error",
	}
	JsonBindingError = ErrorEnums{
		ErrorCode:    -5,
		ErrorMessage: "Json binding error",
	}
	AuthenticateFailure = ErrorEnums{
		ErrorCode:    -6,
		ErrorMessage: "Authenticate fail",
	}
	Unauthorized = ErrorEnums{
		ErrorCode:    -7,
		ErrorMessage: "Unauthorized",
	}
	DataFormatError = ErrorEnums{
		ErrorCode:    -8,
		ErrorMessage: "Data format error",
	}
	Forbidden = ErrorEnums{
		ErrorCode:    -9,
		ErrorMessage: "You don't have permission to perform this action",
	}
	UserAccountAlreadyActived = ErrorEnums{
		ErrorCode:    1001,
		ErrorMessage: "User account already actived",
	}
	UserAlreadyConfigureOtp = ErrorEnums{
		ErrorCode:    1002,
		ErrorMessage: "User already configure otp",
	}
	CreateFolderError = ErrorEnums{
		ErrorCode:    1003,
		ErrorMessage: "Can not create folder",
	}
	ListFolderError = ErrorEnums{
		ErrorCode:    1004,
		ErrorMessage: "Can not list folder",
	}
)
