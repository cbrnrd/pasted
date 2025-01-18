package transforms

import "io"


// Transformer defines an interface for bi-directional transformations on an io.Reader.
type Transformer interface {
	Transform(input io.Reader) (io.Reader, error)
	ReverseTransform(input io.Reader) (io.Reader, error)
}
