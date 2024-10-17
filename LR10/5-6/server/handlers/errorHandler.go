package handlers

import "net/http"

func HandleError(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
}
