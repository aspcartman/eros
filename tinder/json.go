package tinder

import (
	"io"
	"encoding/json"
	"github.com/aspcartman/exceptions"
	"bytes"
)

type jsonm map[string]interface{}

type breadcloser struct {
	io.Reader
}

func (breadcloser) Close() error {
	return nil
}

func reader(j interface{}) io.ReadCloser {
	if j == nil {
		return nil
	}
	data, err := json.Marshal(j)
	if err != nil {
		e.Throw("failed marshalling request data", err, e.Map{
			"data": j,
		})
	}

	return breadcloser{bytes.NewReader(data)}
}

func btom(b []byte) jsonm {
	m := jsonm{}
	e.Must(json.Unmarshal(b, &m), "failed unmarshalling tinder response to map", e.Map{
		"response": string(b),
	})
	return m
}
