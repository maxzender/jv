package formatters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func Indent(r io.Reader) (io.Reader, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, fmt.Errorf("parse error: %v", err)
	}

	content, err = json.MarshalIndent(data, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	return bytes.NewReader(content), nil
}
