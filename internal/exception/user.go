package exception

import "net/http"

var (
	ExceptionUserNotFound           = New(http.StatusNotFound, 2000, "User not found")
	ExceptionUserEmailAlreadyExists = New(http.StatusBadRequest, 2001, "User email already exists")
	ExceptionUserPasswordIncorrect  = New(http.StatusBadRequest, 2002, "User password incorrect")
)
