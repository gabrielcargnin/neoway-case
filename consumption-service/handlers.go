package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"neoway-case/db"
	"neoway-case/errors"
	"neoway-case/util"
	"net/http"
)

type Payload struct {
	Errors  []errors.ResponseError `json:"errors,omitempty"`
	Message string                 `json:"message"`
	Code    int                    `json:"code"`
}

func insertConsumptionHandler(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		renderError(w, errors.E(errors.Message("Could not read file"), http.StatusInternalServerError, errors.Op("main.uploadFileHandler")))
	}
	reader := bytes.NewReader(buf)
	err = insertConsumption(r, reader)
	if err != nil {
		errors.LogCleanStackTrace(err)
		renderError(w, errors.E(errors.Op("main.uploadFileHandler"), err, errors.Message("Could not process upload request")))
		return
	}
	renderOk(w)
}

func insertConsumption(r *http.Request, reader io.Reader) error {
	const errorMessage errors.Message = "Error trying to insert data"
	const op errors.Op = "main.insertConsumption"
	rows, err := util.Parse(reader)
	if err != nil {
		return errors.E(op, errorMessage, err)
	}
	if err := util.Validate(rows); err != nil {
		return errors.E(op, errorMessage, err)
	}
	if err := db.InsertConsumption(r.Context(), rows); err != nil {
		return errors.E(err, errors.Message("Consumption insert failed"), errors.Op("main.insertConsumption"))
	}
	return nil
}

func renderOk(w http.ResponseWriter) {
	p := Payload{
		Message: "Successfully stored file data",
		Code:    http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

func renderError(w http.ResponseWriter, err error) {
	code := errors.Code(err)
	var res []errors.ResponseError
	p := Payload{
		Errors:  append(res, errors.GetResponseErr(err)),
		Message: "Could not complete request",
		Code:    code,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(p)
}
