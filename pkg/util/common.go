package util

type ErrorInfo struct {
	HttpCode int
	Message  string
	Err      error
}

const (
	ValidationErrorMessage = "Validation error"
	BindingErrorMessage    = "Binding error"
	InternalServiceError   = "InternalServiceError"
	InvalidAccountAddress  = "InvalidAccountAddress"
)
