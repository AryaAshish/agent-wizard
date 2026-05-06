package cli

// ExitCoder allows main to exit with a stable non-zero exit code without abusing sentinel strings.
type ExitCoder interface {
	error
	Code() int
}

type exitCode struct {
	code int
	err  error
}

func (e *exitCode) Error() string { return e.err.Error() }
func (e *exitCode) Code() int     { return e.code }

// NewExit wraps an error with a specific OS exit code.
func NewExit(code int, err error) error {
	if err == nil {
		return nil
	}
	return &exitCode{code: code, err: err}
}
