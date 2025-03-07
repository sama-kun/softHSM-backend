package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
)

// DecodeJSON декодирует JSON из запроса.
func DecodeJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}

	return nil
}
