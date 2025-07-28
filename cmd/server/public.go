package main

import (
	"encoding/json"
	"go-mini-server/internal/errors"
	"io"
	"net/http"
)

func handleUserGet(r *http.Request) (interface{}, error) {
	body := struct {
		Id int `json:"id"`
	}{}

	err := decodeBody(r, &body)
	if err != nil {
		return nil, errors.InvalidRequest
	}

	u, err := userService.FetchById(body.Id)
	if err != nil {
		return nil, err
	}

	return u, err
}

func decodeBody(r *http.Request, v interface{}) error {
	if r.Method != http.MethodPost {
		return errors.InvalidRequest
	}

	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &v); err != nil {
		return errors.InvalidRequest
	}

	return nil
}
