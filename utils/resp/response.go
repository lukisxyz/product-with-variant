package resp

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, msg string, status int, data, meta any) error {
	var j struct {
		Msg  string `json:"msg"`
		Data any    `json:"data"`
		Meta any    `json:"meta"`
	}

	j.Data = data
	j.Meta = meta
	j.Msg = msg

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(j)
}

func writeMessage(w http.ResponseWriter, status int, msg string) error {
	var j struct {
		Msg string `json:"msg"`
	}

	j.Msg = msg
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(j)
}

func WriteError(w http.ResponseWriter, status int, err error) error {
	return writeMessage(w, status, err.Error())
}
