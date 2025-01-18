package transforms

import "io"

// ChainTransformer chains multiple Transformers together, supporting both forward and reverse transformations.
type ChainTransformer struct {
	transformers []Transformer
}

// NewChainTransformer creates a new ChainTransformer.
func NewChainTransformer(transformers ...Transformer) *ChainTransformer {
	return &ChainTransformer{transformers: transformers}
}

// Transform applies all transformers in sequence.
func (ct *ChainTransformer) Transform(input io.Reader) (io.Reader, error) {
	var err error
	current := input
	for _, transformer := range ct.transformers {
		current, err = transformer.Transform(current)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}

// ReverseTransform applies all transformers in reverse order.
func (ct *ChainTransformer) ReverseTransform(input io.Reader) (io.Reader, error) {
	var err error
	current := input
	for i := len(ct.transformers) - 1; i >= 0; i-- {
		current, err = ct.transformers[i].ReverseTransform(current)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}
