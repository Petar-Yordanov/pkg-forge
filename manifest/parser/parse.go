package parser

import (
	"bytes"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseFile(path string) ([]Document, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(bytes.NewReader(b))
}

func Parse(r io.Reader) ([]Document, error) {
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)

	var docs []Document
	for {
		var d Document
		err := dec.Decode(&d)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(d.Entries) == 0 {
			continue
		}
		docs = append(docs, d)
	}
	return docs, nil
}
