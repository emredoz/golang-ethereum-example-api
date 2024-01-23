package util

type ErrorInfo struct {
	HttpCode int
	Message  string
	Err      error
}

const (
	ValidationErrorMessage            = "validation error"
	BindingErrorMessage               = "binding error"
	InternalServiceErrorMessage       = "internal Service Error"
	InvalidAccountAddressErrorMessage = "invalid Account Address"
)
