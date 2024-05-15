package main

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func JSONReponse(w http.ResponseWriter, code int, message string, err error) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	errResp := ErrorResponse{Message: message}
	if err != nil {
		errResp.Error = err.Error()
	}
	err = json.NewEncoder(w).Encode(errResp)
	return err
}

func ServerErrorResp(w http.ResponseWriter, message string, err error) error {
	return JSONReponse(w, http.StatusInternalServerError, message, err)
}

func BadRequestResp(w http.ResponseWriter, message string, err error) error {
	return JSONReponse(w, http.StatusBadRequest, message, err)
}

func NotFoundResp(w http.ResponseWriter, message string, err error) error {
	return JSONReponse(w, http.StatusNotFound, message, err)
}

func NotAllowed(w http.ResponseWriter, message string, err error) error {
	return JSONReponse(w, http.StatusMethodNotAllowed, message, err)
}
