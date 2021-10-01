package formatter

import (
	"encoding/json"

	"github.com/profe-ajedrez/sapo/pkg/morestring"
)

type Jsonificable interface {
	ToJson() (string, error)
	ToPrettyJson() (string, error)
}

func ToJson(elem interface{}) (string, error) {
	json, err := json.Marshal(elem)
	if err == nil {
		return string(json), err
	}

	return morestring.Null, err
}

func ToPrettyJson(elem interface{}) (string, error) {
	json, err := json.MarshalIndent(elem, "   ", "   ")
	if err == nil {
		return string(json), err
	}

	return morestring.Null, err
}
