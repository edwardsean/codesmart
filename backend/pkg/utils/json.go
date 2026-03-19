package utils

import (
	"encoding/json"
	"errors"
	"io"
)

func ParseJson(body io.Reader, payload any) error {
	if body == nil {
		return errors.New("missing body for parsing")
	}

	return json.NewDecoder(body).Decode(payload)
}
