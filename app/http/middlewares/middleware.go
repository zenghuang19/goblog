package middlewares

import "net/http"

// HttpHandlerFunc įŽå ââ func(http.ResponseWriter, *http.Request)
type HttpHandlerFunc func(w http.ResponseWriter, r *http.Request)
