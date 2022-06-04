package encoding

import (
	"io"

	"gopkg.in/yaml.v3"
)

type YamlDecoder struct{}

func (yd YamlDecoder) Decode(r io.Reader, v any) error {
	return yaml.NewDecoder(r).Decode(v)
}
