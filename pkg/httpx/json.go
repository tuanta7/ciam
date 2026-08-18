package httpx

import (
	"encoding/json"
	"io"
	"net/http"
)

type JSON map[string]any

func ReadJSON(payload io.Reader, data any) error {
	decoder := json.NewDecoder(payload)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func ReadAndValidateJSON(payload io.Reader, data any) error {
	err := ReadJSON(payload, data)
	if err != nil {
		return err
	}

	return ValidateStruct(data)
}

func WriteJSON(w http.ResponseWriter, code int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = w.Write(jsonData)
	return err
}

func ErrorJSON(w http.ResponseWriter, err Error) error {
	return WriteJSON(w, err.Code, map[string]string{
		"error": err.Error(),
		"hint":  err.Hint,
	})
}
