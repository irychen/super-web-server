package exception

import "fmt"

type Exception struct {
	StatusCode int      `json:"-"`
	Code       int      `json:"code"`
	Message    string   `json:"message"`
	Details    []string `json:"details"`
}

var codes = map[int]string{}

func New(statusCode int, code int, message string) *Exception {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("Exception code %d already exists", code))
	}
	codes[code] = message
	return &Exception{StatusCode: statusCode, Code: code, Message: message, Details: []string{}}
}

func (e *Exception) clone() *Exception {
	ne := *e
	return &ne
}

func (e *Exception) AppendDetails(details ...string) *Exception {
	e = e.clone()
	e.Details = append(e.Details, details...)
	return e
}

func (e *Exception) Is(err *Exception) bool {
	if e == nil || err == nil {
		return false
	}
	return e.Code == err.Code
}
