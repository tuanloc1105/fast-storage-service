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
const BytesInMB int32 = 1024 * 1024
const BytesInKB int32 = 1024
const LogFileLocation = "fast_storage_service_log_%d_%d_%d.log"
const DeltaPositive = 0.5
const DeltaNegative = -0.5
const YyyyMmDdHhMmSsFormat = "2006-01-02 15:04:05"
const FileStatDateTimeLayout = "2006-01-02 15:04:05.999999999 -0700"
const RarFileTimeLayout = "20060102150405"
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
		ErrorMessage: "Cannot create folder",
	}
	ListFolderError = ErrorEnums{
		ErrorCode:    1004,
		ErrorMessage: "Cannot list folder",
	}
	EmptyFileInformationError = ErrorEnums{
		ErrorCode:    1005,
		ErrorMessage: "Empty file information",
	}
	CountFileError = ErrorEnums{
		ErrorCode:    1006,
		ErrorMessage: "Cannot count file",
	}
	SaveFileError = ErrorEnums{
		ErrorCode:    1007,
		ErrorMessage: "Cannot save file",
	}
	DownloadFileError = ErrorEnums{
		ErrorCode:    1008,
		ErrorMessage: "Download file error",
	}
	FileStatictisError = ErrorEnums{
		ErrorCode:    1009,
		ErrorMessage: "File statictis error",
	}
	CheckMaximunStorageError = ErrorEnums{
		ErrorCode:    1010,
		ErrorMessage: "Check maximun storage error",
	}
	UploadFileSizeExceeds = ErrorEnums{
		ErrorCode:    1011,
		ErrorMessage: "Upload file size exceeds the limit",
	}
	RemoveFileError = ErrorEnums{
		ErrorCode:    1012,
		ErrorMessage: "Cannot remove file",
	}
	OtpError = ErrorEnums{
		ErrorCode:    1013,
		ErrorMessage: "Otp error",
	}
	WrongOtpError = ErrorEnums{
		ErrorCode:    1014,
		ErrorMessage: "Wrong otp",
	}
	FolderNotExistError = ErrorEnums{
		ErrorCode:    1015,
		ErrorMessage: "Folder does not exist",
	}
	FolderAlreadySecureError = ErrorEnums{
		ErrorCode:    1016,
		ErrorMessage: "Folder already secured",
	}
	InputOtpEmptyError = ErrorEnums{
		ErrorCode:    1017,
		ErrorMessage: "Input otp empty",
	}
	SecureFolderInvalidCredentialError = ErrorEnums{
		ErrorCode:    1018,
		ErrorMessage: "Cannot view secure folder due to wrong `PASSWORD` or `OTP`",
	}
	CreateFileError = ErrorEnums{
		ErrorCode:    1019,
		ErrorMessage: "Cannot create file",
	}
	RenameFolderError = ErrorEnums{
		ErrorCode:    1020,
		ErrorMessage: "Cannot rename folder",
	}
	RenameNonexistentDirectoryError = ErrorEnums{
		ErrorCode:    1021,
		ErrorMessage: "Rename nonexistent folder error",
	}
	RenameSecuredDirectoryError = ErrorEnums{
		ErrorCode:    1021,
		ErrorMessage: "Cannot rename a secured folder",
	}
	HashPasswordForSecuredFolderError = ErrorEnums{
		ErrorCode:    1022,
		ErrorMessage: "Cannot secure folder due to hash password error",
	}
	ZipFolderError = ErrorEnums{
		ErrorCode:    1023,
		ErrorMessage: "Cannot zip folder",
	}
	EmptyUserToShareFileOrFolderError = ErrorEnums{
		ErrorCode:    1024,
		ErrorMessage: "Empty user to share file or folder",
	}
	InvalidListOfUserEmailToShareError = ErrorEnums{
		ErrorCode:    1025,
		ErrorMessage: "Invalid list of user email to share",
	}
)
