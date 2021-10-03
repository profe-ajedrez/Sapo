package errors

const AS_PLANE_STRING = 0
const AS_JSON = 1

func ErrorNoFile() error {
	return New("Config file not exists")
}

func ErrorEmptyFileName() error {
	return New("File path not provided")
}

func ErrorNoTable() error {
	return New("Table name not provided")
}

func ErrorUndetermined() error {
	return New("Undetermined Error")
}

// New returns an error that formats as the given text.
func New(message string) error {
	return &errorString{message}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func ToJson(err error) string {
	return "{\"error\":\"" + err.Error() + "\"}"
}
