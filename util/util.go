package util

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

// ShiftPath Can remove first two chunks of URL path.
// Input: "/api/user/1234/"
// Output: "/user/1234", "/1234"
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	// if no "/" left in path return "/"
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

// AddCORSHeaders Adds header that allow CORS
func AddCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

// CreateHTTPResponse responds with status code, message, and logs the error
func CreateHTTPResponse(w http.ResponseWriter, message string, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
	if err != nil {
		fmt.Println(err)
	}
}
