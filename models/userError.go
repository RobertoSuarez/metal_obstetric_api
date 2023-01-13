package models

type UserError struct {
	Message string `json:"message"`
	IsUser  bool   `json:"isUser"`
}

func (userError UserError) Error() string {
	return userError.Message
}

func newUserError(msg string, isUser bool) error {
	return UserError{
		Message: msg,
		IsUser:  isUser,
	}
}
