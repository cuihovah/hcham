package uc

import (
	"errors"
)

func _set(_type string, userId string, value string) (*UserIndex, error) {
	if value != "" {
		index := &UserIndex{}
		index.UserId = userId
		index.Id = value
		index.Type = _type
		return index, nil
	}
	return nil, errors.New("Value is empty")
}