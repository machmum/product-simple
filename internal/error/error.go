package errin

// Error contains code and error value
// nil value represent no error.
type Error interface {
	error
	code
	shadowed
}

type code interface {
	Code() int
}

type shadowed interface {
	Original() error
	AsMessage() bool
}

// NewError: If debug is false, omit error message and return shadowed error
func NewError(code int, err error) Error {
	return &customErr{c: code, e: err}
}

func NewErrMessage(code int, err error) Error {
	return &customErr{c: code, e: err, eMsg: true}
}

type customErr struct {
	c    int
	e    error
	eMsg bool
}

func (c *customErr) Error() string {
	return c.e.Error()
}

func (c *customErr) Code() int {
	return c.c
}

func (c *customErr) Original() error {
	return c.e
}

func (c *customErr) AsMessage() bool {
	return c.eMsg
}
