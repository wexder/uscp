package encoding

import (
	"encoding/json"
	"io"
)

type JsonDecoder struct{}

func (jd JsonDecoder) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
