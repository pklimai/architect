package business_error

import "errors"

// GetCode returns code of Error.
// If err is not Error, then returns Internal.
func GetCode(err error) Code {
	if err == nil {
		return OK
	}

	if businessError := new(Error); errors.As(err, &businessError) {
		return businessError.GetCode()
	}

	return Internal
}
