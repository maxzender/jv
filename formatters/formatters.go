package formatters

import "io"

type Formatter func(io.Reader) (io.Reader, error)

func Apply(reader io.Reader, funcs ...Formatter) (io.Reader, error) {
	r := reader
	var err error
	for _, formatter := range funcs {
		r, err = formatter(r)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}
