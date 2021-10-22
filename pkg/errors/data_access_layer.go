package errors

import (
	"fmt"
)

type UserDAOError struct {
	Description string
	Errors      uint16
	Err         error
}

/*--------------------------DAO calls errors--------------------------*/
const (
	UserNotFoundInDB uint16 = 1 << iota
	ReceiverNotFoundInDB
	ReceiverNotInList
	SameUser
)

/*-------------------------------------------------------------------*/

var UserDAOErrorDescriptionMap = map[int]string{
	int(UserNotFoundInDB):     "User not found.",
	int(ReceiverNotFoundInDB): "Receiver found.",
	int(ReceiverNotInList):    "Receiver not in list.",
	int(SameUser):             "Receiver can not be the same user.",
}

func (r *UserDAOError) Error() string {
	return fmt.Sprintf("desc %v: err %v", r.Description, r.Err)
}

func NewUserDAOError(errors uint16, err error) *UserDAOError {
	desc := UserDAOErrorDescriptionMap[int(errors)]
	return &UserDAOError{
		Description: desc,
		Errors:      errors,
		Err:         err,
	}
}
