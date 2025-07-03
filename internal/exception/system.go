package exception

import "net/http"

var (
	ExceptionInternalServerError = New(http.StatusInternalServerError, 1000, "Internal server error")
	ExceptionNotFound            = New(http.StatusNotFound, 1001, "Not found")
	ExceptionBadRequest          = New(http.StatusBadRequest, 1002, "Bad request")
	ExceptionInvalidParam        = New(http.StatusBadRequest, 1003, "Invalid param")
	ExceptionUnauthorized        = New(http.StatusUnauthorized, 1004, "Unauthorized")
	ExceptionTokenNotFound       = New(http.StatusUnauthorized, 1005, "Token not found")
	ExceptionTokenExpired        = New(http.StatusUnauthorized, 1006, "Token expired")
	ExceptionTokenGenerateFailed = New(http.StatusUnauthorized, 1007, "Token generate failed")
	ExceptionForbidden           = New(http.StatusForbidden, 1008, "Forbidden")
	ExceptionTooManyRequests     = New(http.StatusTooManyRequests, 1009, "Too many requests")
	ExceptionBadGateway          = New(http.StatusBadGateway, 1010, "Bad gateway")
	ExceptionServiceUnavailable  = New(http.StatusServiceUnavailable, 1011, "Service unavailable")
	ExceptionGatewayTimeout      = New(http.StatusGatewayTimeout, 1012, "Gateway timeout")
	ExceptionNotImplemented      = New(http.StatusNotImplemented, 1013, "Not implemented")
	ExceptionServiceError        = New(http.StatusServiceUnavailable, 1014, "Service error")
	ExceptionServiceTimeout      = New(http.StatusServiceUnavailable, 1015, "Service timeout")
	ExceptionDatabaseError       = New(http.StatusInternalServerError, 1016, "Database error")
)
