package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const yamlDelimiter = `---`

var errFileNotFound = errors.New("Specified file was not found")

type storedFile struct {
	Meta    map[string]interface{} `json:"meta"`
	Content string                 `json:"content"`
}

func loadStoredFile(filename string) (*storedFile, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, errFileNotFound
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read file")
	}

	return storedFileFromString(string(content))
}

func storedFileFromString(content string) (*storedFile, error) {
	// Look at first line and see whether this file has a metadata part
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) == 0 {
		// Empty file
		return &storedFile{}, nil
	}

	var (
		metadata     []string
		contentStart int
	)

	if lines[0] == yamlDelimiter {
		// This file has a metadata part
		for i := 1; i < len(lines); i++ {
			if lines[i] == yamlDelimiter {
				contentStart = i + 1
				break
			}

			metadata = append(metadata, lines[i])
		}
	}

	file := &storedFile{
		Content: strings.TrimSpace(strings.Join(lines[contentStart:], "\n")),
		Meta:    map[string]interface{}{},
	}

	if len(metadata) > 0 {
		if err := yaml.NewDecoder(strings.NewReader(strings.Join(metadata, "\n"))).Decode(&file.Meta); err != nil {
			return nil, errors.Wrap(err, "Unable to parse metadata part")
		}
	}

	return file, nil
}

func (s storedFile) GetMetaString(key string) string {
	if v, ok := s.Meta[key].(string); ok {
		return v
	}
	return ""
}

func (s storedFile) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "Unable to create file")
	}
	defer f.Close()

	if len(s.Meta) > 0 {
		fmt.Fprintln(f, yamlDelimiter)
		if err := yaml.NewEncoder(f).Encode(s.Meta); err != nil {
			return errors.Wrap(err, "Unable to write metadata")
		}
		fmt.Fprintln(f, yamlDelimiter)
	}

	fmt.Fprintln(f, s.Content)

	return nil
}
