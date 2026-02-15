package parse

import (
	"os"

	"github.com/Petar-Yordanov/pkg-forge/manifest"
	"gopkg.in/yaml.v3"
)

func ParseBytes(b []byte) (*manifest.Document, error) {
	var doc manifest.Document
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func ParseString(s string) (*manifest.Document, error) {
	return ParseBytes([]byte(s))
}

func ParseFile(path string) (*manifest.Document, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseBytes(b)
}
