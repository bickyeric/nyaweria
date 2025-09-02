package errors

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field"`
}

func (e ErrorDetail) Error() string {
	return e.Message
}
