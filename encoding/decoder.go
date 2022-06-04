package encoding

import "io"

type Decoder interface {
	Decode(r io.Reader, v any) error
}
