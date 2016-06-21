package server

import "net/http"

//httpError is wrapper around error that holds status code information
type httpError struct {
	error      `json:"-"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

//HTTPErrorStatusCode returns a status code.
func (e httpError) HTTPErrorStatusCode() int {
	return e.StatusCode
}

func (e httpError) Render(w http.ResponseWriter) {
	w.WriteHeader(e.StatusCode)
	w.Write([]byte(e.Message))
}

//NewErrorWithStatusCode allows you to associate
//a specific HTTP Status Code to an error.
//The Server will take that code and set
//it as the response status.
func NewErrorWithStatusCode(err error, code int) *httpError {
	return &httpError{err, code, err.Error()}
}

//NewBadRequestError creates a new API error
//that has the 400 HTTP status code associated to it.
func NewBadRequestError(err error) *httpError {
	return NewErrorWithStatusCode(err, http.StatusBadRequest)
}

//NewForbiddenError create a new APU error
//that has the 403 HTTP status code associated to it.
func NewForbiddenError(err error) *httpError {
	return NewErrorWithStatusCode(err, http.StatusForbidden)
}

//NewRequestNotFoundError creates a new API error
//that has the 404 HTTP status code associated to it.
func NewRequestNotFoundError(err error) *httpError {
	return NewErrorWithStatusCode(err, http.StatusNotFound)
}

//NewRequestConflictError creates a new API error
//that has the 409 HTTP status code associated to it.
func NewRequestConflictError(err error) *httpError {
	return NewErrorWithStatusCode(err, http.StatusConflict)
}

//NewRequesForbiddenError creates a new API error
//that has the 403 HTTP status code associated to it.
func NewRequesForbiddenError(err error) *httpError {
	return NewErrorWithStatusCode(err, http.StatusForbidden)
}

//NewServerError creates a new API error
//that has the 500 HTTP status code associated to it.
func NewServerError(err error) *httpError {
	return NewErrorWithStatusCode(err, http.StatusInternalServerError)
}
