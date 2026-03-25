package verification

const ErrCodeProviderSessionCreateFailed = "PROVIDER_SESSION_CREATE_FAILED"
const ErrCodeRequestNotFound = "REQUEST_NOT_FOUND"

type Error struct {
	Code string
	Err  error
}

func (err *Error) Error() string {
	if err.Err == nil {
		return err.Code
	}

	return err.Code + ": " + err.Err.Error()
}

func (err *Error) Unwrap() error {
	return err.Err
}
