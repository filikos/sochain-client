package sochain

type ClientError struct {
	err        error
	statuscode int
}

func (c ClientError) Error() string {
	return c.err.Error()
}

func (c ClientError) Code() int {
	return c.statuscode
}

func NewClientErr(e error, statuscode int) *ClientError {
	return &ClientError{
		err:        e,
		statuscode: statuscode,
	}
}
